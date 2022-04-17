use crate::constants::*;
use crate::errors::*;
use crate::helper_fns::*;
use crate::ix_accounts::*;
use crate::state::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use anchor_lang::Discriminator;
use anchor_spl::token::{self, Burn, Transfer};
use mpl_token_metadata::state::{MAX_NAME_LENGTH, MAX_SYMBOL_LENGTH, MAX_URI_LENGTH};

use std::str::FromStr;

use anchor_lang::solana_program::{
    program::invoke_signed,
    serialize_utils::{read_pubkey, read_u16},
};

use mpl_token_metadata::instruction::{
    create_master_edition, create_metadata_accounts, update_metadata_accounts,
};

pub mod constants;
pub mod errors;
pub mod helper_fns;
pub mod ix_accounts;
pub mod state;
pub mod structs;

declare_id!("5WwhzMCFSgWYxiuKrbsB9wtg9T49Mm1fD1v2UdhD5oYi");

#[program]
pub mod someplace {

    use super::*;

    pub fn create_market_listing(
        ctx: Context<InitMarketListing>,
        index: u64,
        price: u64,
    ) -> ProgramResult {
        let market_authority = &mut ctx.accounts.market_authority;
        let market_listing = &mut ctx.accounts.market_listing;
        let nft_mint = &mut ctx.accounts.nft_mint;
        let market_listing_token_account = &mut ctx.accounts.market_listing_token_account;
        let seller = &mut ctx.accounts.seller;
        let cpi_accounts = Transfer {
            from: ctx.accounts.seller_nft_token_account.to_account_info(),
            to: market_listing_token_account.to_account_info(),
            authority: seller.to_account_info(),
        };
        let cpi_program = ctx.accounts.token_program.to_account_info();
        let cpi = CpiContext::new(cpi_program, cpi_accounts);
        token::transfer(
            cpi,
            (1 as f64 * 10_usize.pow(nft_mint.decimals as u32) as f64) as u64,
        )?;

        if market_authority.listings != index {
            return Err(QuestError::SuspiciousTreasury.into());
        }
        market_listing.market_authority = market_authority.key();
        market_listing.seller_market_token_account = ctx.accounts.seller_market_token_account.key();
        market_listing.nft_mint = nft_mint.key();
        market_listing.index = index;
        market_listing.price =
            (price as f64 * 10_usize.pow(market_authority.market_decimals as u32) as f64) as u64;
        market_listing.fulfilled = 0;

        market_authority.listings += 1;

        Ok(())
    }
    pub fn fulfill_market_listing(
        ctx: Context<FulfillMarketListing>,
        market_authority_bump: u8,
    ) -> ProgramResult {
        let market_authority = &mut ctx.accounts.market_authority;
        let market_listing = &mut ctx.accounts.market_listing;
        let nft_mint = &mut ctx.accounts.nft_mint;
        let market_listing_token_account = &mut ctx.accounts.market_listing_token_account;
        let market_uid = market_authority.market_uid;
        let oracle_key = ctx.accounts.oracle.key();
        let market_authority_bump_bytes = market_authority_bump.to_le_bytes();
        let seeds = &[
            PREFIX.as_ref(),
            MARKET.as_ref(),
            oracle_key.as_ref(),
            market_uid.as_ref(),
            market_authority_bump_bytes.as_ref(),
        ];
        let market_authority_signer = &[&seeds[..]];

        // send nft mint from market_authority ata to buyer ata
        token::transfer(
            CpiContext::new_with_signer(
                ctx.accounts.token_program.to_account_info(),
                Transfer {
                    from: market_listing_token_account.to_account_info(),
                    to: ctx.accounts.buyer_nft_token_account.to_account_info(),
                    authority: market_authority.to_account_info(),
                },
                market_authority_signer,
            ),
            (1 as f64 * 10_usize.pow(nft_mint.decimals as u32) as f64) as u64,
        )?;

        // tender market mint from buyer ata to seller ata
        let seller_market_token_account = &mut ctx.accounts.seller_market_token_account;
        if market_listing.seller_market_token_account != seller_market_token_account.key() {
            return Err(QuestError::SuspiciousTreasury.into());
        }
        token::transfer(
            CpiContext::new(
                ctx.accounts.token_program.to_account_info(),
                Transfer {
                    from: ctx.accounts.buyer_market_token_account.to_account_info(),
                    to: seller_market_token_account.to_account_info(),
                    authority: ctx.accounts.buyer.to_account_info(),
                },
            ),
            market_listing.price,
        )?;

        market_listing.fulfilled = 1;

        Ok(())
    }
    pub fn create_listing(
        ctx: Context<CreateListing>,
        config_index: u64,
        price: u64,
        lifecycle_start: u64,
    ) -> ProgramResult {
        let listing = &mut ctx.accounts.listing;
        let treasury_authority = &ctx.accounts.treasury_authority;
        listing.treasury_authority = treasury_authority.key();
        listing.batch = ctx.accounts.batch.key();
        listing.oracle = ctx.accounts.oracle.key();
        listing.config_index = config_index;
        listing.price = (price as f64
            * 10_usize.pow(treasury_authority.treasury_decimals as u32) as f64)
            as u64;
        listing.lifecycle_start = lifecycle_start;

        Ok(())
    }

    pub fn enable_batch_uploading(ctx: Context<EnableBatches>) -> ProgramResult {
        let batches = &mut ctx.accounts.batches;
        batches.counter = 0;
        batches.oracle = ctx.accounts.oracle.key().clone();

        Ok(())
    }

    pub fn init_market(
        ctx: Context<InitMarket>,
        market_uid: Pubkey,
        name: String,
    ) -> ProgramResult {
        let market = &mut ctx.accounts.market_authority;
        let market_mint = &mut ctx.accounts.market_mint;
        let oracle = &ctx.accounts.oracle;

        market.market_decimals = market_mint.decimals;
        market.listings = 0;
        market.name = name;
        market.market_mint = market_mint.key();
        market.market_uid = market_uid;
        market.oracle = oracle.key();

        Ok(())
    }

    pub fn init_treasury(ctx: Context<InitTreasury>, adornment: String) -> ProgramResult {
        let authority = &mut ctx.accounts.treasury_authority;
        let treasury_token_account = &mut ctx.accounts.treasury_token_account;
        let treasury_mint = &mut ctx.accounts.treasury_token_mint;

        let oracle = &ctx.accounts.oracle;
        authority.oracle = oracle.key();
        authority.adornment = adornment;
        authority.whitelists = 0;
        authority.treasury_mint = treasury_mint.key();
        authority.treasury_token_account = treasury_token_account.key();

        let cpi_accounts = Transfer {
            from: ctx.accounts.oracle_token_account.to_account_info(),
            to: treasury_token_account.to_account_info(),
            authority: oracle.to_account_info(),
        };
        let cpi_program = ctx.accounts.token_program.to_account_info();
        let cpi = CpiContext::new(cpi_program, cpi_accounts);
        token::transfer(
            cpi,
            (100000 as f64 * 10_usize.pow(treasury_mint.decimals as u32) as f64) as u64,
        )?;

        Ok(())
    }

    pub fn add_whitelisted_cm(
        ctx: Context<AddWhitelistedCM>,
        candy_machine_creator: Pubkey,
        candy_machine: Pubkey,
    ) -> ProgramResult {
        let oracle = &ctx.accounts.oracle;
        let authority = &mut ctx.accounts.treasury_authority;
        let whitelist = &mut ctx.accounts.treasury_whitelist;
        // let whitelist_receipt = &mut ctx.accounts.treasury_whitelist_receipt;
        whitelist.whitelist_id = authority.whitelists;
        whitelist.candy_machine_id = candy_machine;
        whitelist.candy_machine_creator = candy_machine_creator;
        whitelist.treasury_authority = authority.key();
        whitelist.oracle = oracle.key();

        /*
        whitelist_receipt.whitelist_id = authority.whitelists;
        whitelist_receipt.candy_machine_id = candy_machine;
        whitelist_receipt.treasury_authority = authority.key();
        whitelist_receipt.oracle = oracle.key();
        */

        authority.whitelists += 1;
        Ok(())
    }

    pub fn sell_for<'info>(ctx: Context<SellFor>, treasury_bump: u8) -> ProgramResult {
        let treasury_mint = &mut ctx.accounts.treasury_token_mint;
        let treasury_authority = &mut ctx.accounts.treasury_authority;
        let treasury_whitelist = &mut ctx.accounts.treasury_whitelist;
        if treasury_whitelist.treasury_authority.to_string() != treasury_authority.key().to_string()
        {
            return Err(QuestError::SuspiciousTreasury.into());
        }
        /*
        if treasury_authority.treasury_mint.to_string() != treasury_mint.key().to_string() {
            return Err(QuestError::SuspiciousTreasuryMint.into());
        }
        */

        let depo_mint = &mut ctx.accounts.depo_mint;
        let metadata =
            assert_valid_metadata(&ctx.accounts.metadata.to_account_info(), &depo_mint.key())?;
        let creators = metadata.data.creators.unwrap();

        let is_valid = creators
            .iter()
            .find(|creator| {
                creator.address.to_string() == treasury_whitelist.candy_machine_creator.to_string()
            })
            .is_some();
        msg!(
            "{:?} {:?}",
            creators,
            treasury_whitelist.candy_machine_creator.to_string()
        );
        if is_valid == false {
            return Err(QuestError::SuspiciousCandyMachine.into());
        }

        let burn_cpx = Burn {
            mint: depo_mint.to_account_info(),
            to: ctx.accounts.depo_token_account.to_account_info(),
            authority: ctx.accounts.initializer.to_account_info(),
        };
        let burn_cpx_context =
            CpiContext::new(ctx.accounts.token_program.to_account_info(), burn_cpx);
        token::burn(
            burn_cpx_context,
            (1 as f64 * 10_usize.pow(ctx.accounts.depo_mint.decimals as u32) as f64) as u64,
        )?;
        let transfer_cpx = Transfer {
            from: ctx.accounts.treasury_token_account.to_account_info(),
            to: ctx.accounts.initializer_token_account.to_account_info(),
            authority: treasury_authority.to_account_info(),
        };
        let treasury_bump_bytes = treasury_bump.to_le_bytes();

        let oracle_key = ctx.accounts.oracle.key();
        let seeds = &[
            PREFIX.as_ref(),
            BENEFIT_TOKEN.as_ref(),
            oracle_key.as_ref(),
            treasury_bump_bytes.as_ref(),
        ];
        let signer = &[&seeds[..]];
        let cpi = CpiContext::new_with_signer(
            ctx.accounts.token_program.to_account_info(),
            transfer_cpx,
            signer,
        );
        token::transfer(
            cpi,
            (100 as f64 * 10_usize.pow(treasury_mint.decimals as u32) as f64) as u64,
        )?;
        Ok(())
    }

    pub fn add_config_lines(
        ctx: Context<Sync>,
        index: u32,
        config_lines: Vec<ConfigLine>,
    ) -> ProgramResult {
        let candy_machine = &mut ctx.accounts.batch;
        let account = candy_machine.to_account_info();
        let current_count = get_config_count(&account.data.borrow_mut())?;
        let mut data = account.data.borrow_mut();
        let mut fixed_config_lines = vec![];
        // No risk overflow because you literally cant store this many in an account
        // going beyond u32 only happens with the hidden store candies, which dont use this.
        if index > (candy_machine.data.items_available as u32) - 1 {
            return Err(QuestError::IndexGreaterThanLength.into());
        }
        for line in &config_lines {
            let mut array_of_zeroes = vec![];
            while array_of_zeroes.len() < MAX_NAME_LENGTH - line.name.len() {
                array_of_zeroes.push(0u8);
            }
            let name = line.name.clone() + std::str::from_utf8(&array_of_zeroes).unwrap();

            let mut array_of_zeroes = vec![];
            while array_of_zeroes.len() < MAX_URI_LENGTH - line.uri.len() {
                array_of_zeroes.push(0u8);
            }
            let uri = line.uri.clone() + std::str::from_utf8(&array_of_zeroes).unwrap();
            fixed_config_lines.push(ConfigLine { name, uri })
        }

        let as_vec = fixed_config_lines.try_to_vec()?;
        // remove unneeded u32 because we're just gonna edit the u32 at the front
        let serialized: &[u8] = &as_vec.as_slice()[4..];

        let position = CONFIG_ARRAY_START + 4 + (index as usize) * CONFIG_LINE_SIZE;

        let array_slice: &mut [u8] =
            &mut data[position..position + fixed_config_lines.len() * CONFIG_LINE_SIZE];

        array_slice.copy_from_slice(serialized);

        let bit_mask_vec_start = CONFIG_ARRAY_START
            + 4
            + (candy_machine.data.items_available as usize) * CONFIG_LINE_SIZE
            + 4;

        let mut new_count = current_count;
        for i in 0..fixed_config_lines.len() {
            let position = (index as usize)
                .checked_add(i)
                .ok_or(QuestError::NumericalOverflowError)?;
            let my_position_in_vec = bit_mask_vec_start
                + position
                    .checked_div(8)
                    .ok_or(QuestError::NumericalOverflowError)?;
            let position_from_right = 7 - position
                .checked_rem(8)
                .ok_or(QuestError::NumericalOverflowError)?;
            let mask = u8::pow(2, position_from_right as u32);

            let old_value_in_vec = data[my_position_in_vec];
            data[my_position_in_vec] = data[my_position_in_vec] | mask;
            msg!(
                "My position in vec is {} my mask is going to be {}, the old value is {}",
                position,
                mask,
                old_value_in_vec
            );
            msg!(
                "My new value is {} and my position from right is {}",
                data[my_position_in_vec],
                position_from_right
            );
            if old_value_in_vec != data[my_position_in_vec] {
                msg!("Increasing count");
                new_count = new_count
                    .checked_add(1)
                    .ok_or(QuestError::NumericalOverflowError)?;
            }
        }

        // plug in new count.
        data[CONFIG_ARRAY_START..CONFIG_ARRAY_START + 4]
            .copy_from_slice(&(new_count as u32).to_le_bytes());

        Ok(())
    }

    pub fn initialize_candy_machine(
        ctx: Context<NewBatch>,
        data: CandyMachineData,
        name: String,
    ) -> ProgramResult {
        let batch_receipt = &mut ctx.accounts.batch_receipt;
        let batches = &mut ctx.accounts.batches;
        let candy_machine_account = &mut ctx.accounts.batch_account;
        let collection_name = &name;

        batch_receipt.id = batches.counter.clone();
        batch_receipt.batch_account = candy_machine_account.key();
        batch_receipt.oracle = ctx.accounts.oracle.key();
        batch_receipt.name = collection_name.to_string();
        batches.counter += 1;

        if data.uuid.len() != 6 {
            return Err(QuestError::UuidMustBeExactly6Length.into());
        }

        let mut candy_machine = Batch {
            name: collection_name.to_string(),
            oracle: ctx.accounts.oracle.key(),
            data,
        };

        let mut array_of_zeroes = vec![];
        while array_of_zeroes.len() < MAX_SYMBOL_LENGTH - candy_machine.data.symbol.len() {
            array_of_zeroes.push(0u8);
        }
        let new_symbol =
            candy_machine.data.symbol.clone() + std::str::from_utf8(&array_of_zeroes).unwrap();
        candy_machine.data.symbol = new_symbol;

        let mut new_data = Batch::discriminator().try_to_vec().unwrap();
        new_data.append(&mut candy_machine.try_to_vec().unwrap());
        let mut data = candy_machine_account.data.borrow_mut();
        // god forgive me couldnt think of better way to deal with this
        for i in 0..new_data.len() {
            data[i] = new_data[i];
        }

        let vec_start = CONFIG_ARRAY_START
            + 4
            + (candy_machine.data.items_available as usize) * CONFIG_LINE_SIZE;
        let as_bytes = (candy_machine
            .data
            .items_available
            .checked_div(8)
            .ok_or(QuestError::NumericalOverflowError)? as u32)
            .to_le_bytes();
        for i in 0..4 {
            data[vec_start + i] = as_bytes[i]
        }

        Ok(())
    }

    pub fn mint_nft<'info>(
        ctx: Context<'_, '_, '_, 'info, MintNFT<'info>>,
        creator_bump: u8,
        config_index: u64,
    ) -> ProgramResult {
        let listing = &mut ctx.accounts.listing;
        let candy_machine = &mut ctx.accounts.candy_machine;
        let candy_machine_creator = &ctx.accounts.candy_machine_creator;
        let instruction_sysvar_account = &ctx.accounts.instruction_sysvar_account;

        // ensures valid listing for item in batch
        if listing.oracle != ctx.accounts.oracle.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if listing.batch != candy_machine.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if listing.config_index != config_index {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        if listing.lifecycle_start > Clock::get()?.unix_timestamp as u64 {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        // assert_valid_go_live(payer, clock, candy_machine)?;

        let config_line = get_config_line(&candy_machine, config_index as usize)?;

        let cm_key = candy_machine.key();
        let authority_seeds = [PREFIX, cm_key.as_ref(), &[creator_bump]];

        let mut creators: Vec<mpl_token_metadata::state::Creator> =
            vec![mpl_token_metadata::state::Creator {
                address: candy_machine_creator.key(),
                verified: true,
                share: 0,
            }];

        for c in &candy_machine.data.creators {
            creators.push(mpl_token_metadata::state::Creator {
                address: c.address,
                verified: false,
                share: c.share,
            });
        }

        let metadata_infos = vec![
            ctx.accounts.metadata.to_account_info(),
            ctx.accounts.mint.to_account_info(),
            ctx.accounts.mint_authority.to_account_info(),
            ctx.accounts.payer.to_account_info(),
            ctx.accounts.token_metadata_program.to_account_info(),
            ctx.accounts.token_program.to_account_info(),
            ctx.accounts.system_program.to_account_info(),
            ctx.accounts.rent.to_account_info(),
            candy_machine_creator.to_account_info(),
        ];

        let master_edition_infos = vec![
            ctx.accounts.master_edition.to_account_info(),
            ctx.accounts.mint.to_account_info(),
            ctx.accounts.mint_authority.to_account_info(),
            ctx.accounts.payer.to_account_info(),
            ctx.accounts.metadata.to_account_info(),
            ctx.accounts.token_metadata_program.to_account_info(),
            ctx.accounts.token_program.to_account_info(),
            ctx.accounts.system_program.to_account_info(),
            ctx.accounts.rent.to_account_info(),
            candy_machine_creator.to_account_info(),
        ];

        invoke_signed(
            &create_metadata_accounts(
                *ctx.accounts.token_metadata_program.key,
                *ctx.accounts.metadata.key,
                *ctx.accounts.mint.key,
                *ctx.accounts.mint_authority.key,
                *ctx.accounts.payer.key,
                candy_machine_creator.key(),
                config_line.name,
                candy_machine.data.symbol.clone(),
                config_line.uri,
                Some(creators),
                candy_machine.data.seller_fee_basis_points,
                true,
                candy_machine.data.is_mutable,
            ),
            metadata_infos.as_slice(),
            &[&authority_seeds],
        )?;

        invoke_signed(
            &create_master_edition(
                *ctx.accounts.token_metadata_program.key,
                *ctx.accounts.master_edition.key,
                *ctx.accounts.mint.key,
                candy_machine_creator.key(),
                *ctx.accounts.mint_authority.key,
                *ctx.accounts.metadata.key,
                *ctx.accounts.payer.key,
                Some(candy_machine.data.max_supply),
            ),
            master_edition_infos.as_slice(),
            &[&authority_seeds],
        )?;

        let mut new_update_authority = Some(candy_machine.oracle);

        if !candy_machine.data.retain_authority {
            new_update_authority = Some(ctx.accounts.update_authority.key());
        }

        invoke_signed(
            &update_metadata_accounts(
                *ctx.accounts.token_metadata_program.key,
                *ctx.accounts.metadata.key,
                candy_machine_creator.key(),
                new_update_authority,
                None,
                Some(true),
            ),
            &[
                ctx.accounts.token_metadata_program.to_account_info(),
                ctx.accounts.metadata.to_account_info(),
                candy_machine_creator.to_account_info(),
            ],
            &[&authority_seeds],
        )?;

        let instruction_sysvar_account_info = instruction_sysvar_account.to_account_info();

        let instruction_sysvar = instruction_sysvar_account_info.data.borrow();

        let mut idx = 0;
        let num_instructions =
            read_u16(&mut idx, &instruction_sysvar).map_err(|_| QuestError::InvalidAccountData)?;

        let associated_token =
            Pubkey::from_str("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL").unwrap();

        for index in 0..num_instructions {
            let mut current = 2 + (index * 2) as usize;
            let start = read_u16(&mut current, &instruction_sysvar).unwrap();

            current = start as usize;
            let num_accounts = read_u16(&mut current, &instruction_sysvar).unwrap();
            current += (num_accounts as usize) * (1 + 32);
            let program_id = read_pubkey(&mut current, &instruction_sysvar).unwrap();

            if program_id != someplace::id()
                && program_id != spl_token::id()
                && program_id != anchor_lang::solana_program::system_program::ID
                && program_id != associated_token
            {
                return Err(QuestError::SuspiciousTransaction.into());
            }
        }

        let tender_cpx = Transfer {
            from: ctx.accounts.initializer_token_account.to_account_info(),
            to: ctx.accounts.treasury_token_account.to_account_info(),
            authority: ctx.accounts.payer.to_account_info(),
        };
        let tender_cpx_context =
            CpiContext::new(ctx.accounts.token_program.to_account_info(), tender_cpx);
        token::transfer(tender_cpx_context, listing.price)?;

        Ok(())
    }
}
