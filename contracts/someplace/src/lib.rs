use crate::constants::*;
use crate::errors::*;
use crate::helper_fns::*;
use crate::ix_accounts::*;
use crate::state::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use anchor_lang::Discriminator;
use anchor_spl::token::{self, Burn, MintTo, Transfer};
use mpl_token_metadata::state::{MAX_NAME_LENGTH, MAX_SYMBOL_LENGTH, MAX_URI_LENGTH};
use std::result::Result;

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

declare_id!("GXFE4Ym1vxhbXLBx2RxqL5y1Ee3XyFUqDksD7tYjAi8z");

#[program]
pub mod someplace {

    use super::*;

    pub fn create_market_listing(
        ctx: Context<InitMarketListing>,
        index: u64,
        price: u64,
    ) -> Result<(), Error> {
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
        market_listing.seller = seller.key();
        market_listing.seller_market_token_account = ctx.accounts.seller_market_token_account.key();
        market_listing.nft_mint = nft_mint.key();
        market_listing.index = index;
        market_listing.price = price;
        market_listing.fulfilled = 0;
        market_listing.listed_at = Clock::get()?.unix_timestamp as u64;

        market_authority.listings += 1;

        Ok(())
    }
    pub fn fulfill_market_listing(
        ctx: Context<FulfillMarketListing>,
        market_authority_bump: u8,
    ) -> Result<(), Error> {
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

        market_listing.fulfilled = Clock::get()?.unix_timestamp;

        Ok(())
    }
    pub fn unlist_market_listing(
        ctx: Context<UnlistMarketListing>,
        market_authority_bump: u8,
    ) -> Result<(), Error> {
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
        if market_listing.seller != ctx.accounts.seller.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if market_listing.nft_mint != nft_mint.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        // send nft mint from market_authority ata to buyer ata
        token::transfer(
            CpiContext::new_with_signer(
                ctx.accounts.token_program.to_account_info(),
                Transfer {
                    from: market_listing_token_account.to_account_info(),
                    to: ctx.accounts.seller_nft_token_account.to_account_info(),
                    authority: market_authority.to_account_info(),
                },
                market_authority_signer,
            ),
            (1 as f64 * 10_usize.pow(nft_mint.decimals as u32) as f64) as u64,
        )?;

        market_listing.fulfilled = -1;

        Ok(())
    }
    pub fn create_listing(
        ctx: Context<CreateListing>,
        config_index: u64,
        price: u64,
        lifecycle_start: u64,
    ) -> Result<(), Error> {
        let listing = &mut ctx.accounts.listing;
        let treasury_authority = &ctx.accounts.treasury_authority;
        listing.treasury_authority = treasury_authority.key();
        listing.batch = ctx.accounts.batch.key();
        listing.oracle = ctx.accounts.oracle.key();
        listing.config_index = config_index;
        listing.price = price;
        listing.lifecycle_start = lifecycle_start;
        listing.is_listed = true;
        listing.mints = 0;

        Ok(())
    }

    pub fn enable_vias(ctx: Context<EnableVias>) -> Result<(), Error> {
        let oracle = &ctx.accounts.oracle;
        let vias = &mut ctx.accounts.vias;
        let treasury_authority = &ctx.accounts.treasury_authority;

        if oracle.key() != treasury_authority.oracle {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        vias.oracle = oracle.key();
        vias.treasury_authority = treasury_authority.key();
        vias.vias = 0;

        Ok(())
    }

    pub fn enable_via_rarity_token_minting(
        ctx: Context<EnableViaRarityTokenMinting>,
        rarity: String,
    ) -> Result<(), Error> {
        let oracle = &ctx.accounts.oracle;
        let via = &mut ctx.accounts.via;
        let via_mapping = &mut ctx.accounts.via_mapping;
        let vias = &mut ctx.accounts.vias;
        let treasury_authority = &ctx.accounts.treasury_authority;
        let rarity_token_mint = &ctx.accounts.rarity_token_mint;

        if oracle.key() != treasury_authority.oracle {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if treasury_authority.key() != vias.treasury_authority {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        via.oracle = oracle.key();
        via.treasury_authority = treasury_authority.key();
        via.token_mint = rarity_token_mint.key();
        via.mints = 0;
        via.rarity = rarity;

        via_mapping.token_mint = rarity_token_mint.key();
        via_mapping.vias_index = vias.vias;

        vias.vias += 1;

        Ok(())
    }

    pub fn modify_listing(
        ctx: Context<ModifyListing>,
        _config_index: u64,
        is_listed: Option<bool>,
        lifecycle_start: Option<u64>,
        price: Option<u64>,
    ) -> Result<(), Error> {
        let oracle = &ctx.accounts.oracle;
        let listing = &mut ctx.accounts.listing;
        let treasury_authority = &ctx.accounts.treasury_authority;

        if oracle.key() != treasury_authority.oracle {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        listing.treasury_authority = treasury_authority.key();
        listing.batch = ctx.accounts.batch.key();
        listing.oracle = ctx.accounts.oracle.key();

        if is_listed.is_some() {
            listing.is_listed = is_listed.unwrap();
        }

        if lifecycle_start.is_some() {
            listing.lifecycle_start = lifecycle_start.unwrap();
        }
        if price.is_some() {
            listing.price = price.unwrap();
        }

        Ok(())
    }

    pub fn enable_batch_uploading(ctx: Context<EnableBatches>) -> Result<(), Error> {
        let batches = &mut ctx.accounts.batches;
        batches.counter = 0;
        batches.oracle = ctx.accounts.oracle.key().clone();

        Ok(())
    }

    pub fn init_market(
        ctx: Context<InitMarket>,
        market_uid: Pubkey,
        name: String,
    ) -> Result<(), Error> {
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

    pub fn init_treasury(ctx: Context<InitTreasury>, adornment: String) -> Result<(), Error> {
        let authority = &mut ctx.accounts.treasury_authority;
        let treasury_token_account = &mut ctx.accounts.treasury_token_account;
        let treasury_mint = &mut ctx.accounts.treasury_token_mint;

        let oracle = &ctx.accounts.oracle;
        authority.oracle = oracle.key();
        authority.adornment = adornment;
        authority.whitelists = 0;
        authority.treasury_decimals = treasury_mint.decimals;
        authority.treasury_mint = treasury_mint.key();
        authority.treasury_token_account = treasury_token_account.key();

        Ok(())
    }

    pub fn add_whitelisted_cm(
        ctx: Context<AddWhitelistedCM>,
        candy_machine_creator: Pubkey,
        candy_machine: Pubkey,
    ) -> Result<(), Error> {
        let oracle = &ctx.accounts.oracle;
        let authority = &mut ctx.accounts.treasury_authority;
        let whitelist = &mut ctx.accounts.treasury_whitelist;
        whitelist.whitelist_id = authority.whitelists;
        whitelist.candy_machine_id = candy_machine;
        whitelist.candy_machine_creator = candy_machine_creator;
        whitelist.treasury_authority = authority.key();
        whitelist.oracle = oracle.key();

        authority.whitelists += 1;
        Ok(())
    }

    pub fn ammend_storefront_splits(
        ctx: Context<AmmendStorefrontSplits>,
        storefront_splits: Vec<Split>,
    ) -> Result<(), Error> {
        let oracle = &ctx.accounts.oracle;
        let authority = &mut ctx.accounts.treasury_authority;
        if authority.oracle != oracle.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        authority.splits = storefront_splits.clone();
        if storefront_splits.len() == 0 {
            return Ok(());
        }

        let sum_shares = authority
            .splits
            .iter()
            .map(|split| split.share)
            .reduce(|accumulator, split| accumulator + split);
        if !sum_shares.is_some() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        let sum = sum_shares.unwrap();
        if sum != 100 {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        Ok(())
    }

    pub fn sell_for<'info>(ctx: Context<SellFor>, treasury_bump: u8) -> Result<(), Error> {
        let treasury_mint = &mut ctx.accounts.treasury_token_mint;
        let treasury_authority = &mut ctx.accounts.treasury_authority;
        let treasury_whitelist = &mut ctx.accounts.treasury_whitelist;
        if treasury_whitelist.treasury_authority.to_string() != treasury_authority.key().to_string()
        {
            return Err(QuestError::SuspiciousTreasury.into());
        }

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
            from: ctx.accounts.depo_token_account.to_account_info(),
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
    ) -> Result<(), Error> {
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
            while array_of_zeroes.len() < MAX_CARDINALITY_LENGTH - line.cardinality.len() {
                array_of_zeroes.push(0u8);
            }
            let cardinality =
                line.cardinality.clone() + std::str::from_utf8(&array_of_zeroes).unwrap();

            let mut array_of_zeroes = vec![];
            while array_of_zeroes.len() < MAX_URI_LENGTH - line.uri.len() {
                array_of_zeroes.push(0u8);
            }
            let uri = line.uri.clone() + std::str::from_utf8(&array_of_zeroes).unwrap();
            fixed_config_lines.push(ConfigLine {
                name,
                cardinality,
                uri,
            })
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
    ) -> Result<(), Error> {
        let batch_receipt = &mut ctx.accounts.batch_receipt;
        let batches = &mut ctx.accounts.batches;
        let candy_machine_account = &mut ctx.accounts.batch_account;
        let collection_name = &name;

        batch_receipt.id = batches.counter.clone();
        batch_receipt.batch_account = candy_machine_account.key();
        batch_receipt.oracle = ctx.accounts.oracle.key();
        batch_receipt.name = collection_name.to_string();
        batch_receipt.items = data.items_available;
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
        ctx: Context<'_, '_, '_, 'info, MintNFTListing<'info>>,
        creator_bump: u8,
        config_index: u64,
    ) -> Result<(), Error> {
        let listing = &mut ctx.accounts.listing;
        let mint_hash = &mut ctx.accounts.mint_hash;
        let candy_machine = &mut ctx.accounts.candy_machine;
        let candy_machine_creator = &ctx.accounts.candy_machine_creator;
        let instruction_sysvar_account = &ctx.accounts.instruction_sysvar_account;

        mint_hash.mint = ctx.accounts.mint.key();
        mint_hash.minter = ctx.accounts.payer.key();
        mint_hash.mint_index = listing.mints;
        mint_hash.fulfilled = Clock::get()?.unix_timestamp;

        listing.mints += 1;

        let treasury_authority = &mut ctx.accounts.treasury_authority;
        if treasury_authority.oracle != ctx.accounts.oracle.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }

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

        if !listing.is_listed {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        // assert_valid_go_live(payer, clock, candy_machine)?;
        let payer = &mut ctx.accounts.payer;
        let mint = &mut ctx.accounts.mint;

        token::mint_to(
            CpiContext::new(
                ctx.accounts.token_program.to_account_info(),
                MintTo {
                    mint: mint.to_account_info(),
                    to: ctx.accounts.mint_ata.to_account_info(),
                    authority: payer.to_account_info(),
                },
            ),
            1,
        )?;

        // r/shittyprogramming - messy code naming sorry :(
        let storefront_splits = treasury_authority.splits.clone();
        let mut i: usize = 0;
        /*
            remaining accounts will be ordered such that -

            the first _n_ accounts must equal to their _n_ account in treasury.splits
            the remainder of remaining accounts present additional batches/candy machines
        */
        if treasury_authority.splits.len() > 0 {
            while i < treasury_authority.splits.len() {
                let split = &storefront_splits[i];
                // the order of token address inclusion is respectful to
                // the sequence in `treasury_authority.splits`
                if ctx.remaining_accounts[i].key() != split.token_address {
                    return Err(QuestError::SuspiciousTransaction.into());
                }
                if split.op_code == 0 {
                    token::burn(
                        CpiContext::new(
                            ctx.accounts.token_program.to_account_info(),
                            Burn {
                                from: ctx.accounts.initializer_token_account.to_account_info(),
                                mint: ctx.remaining_accounts[i].to_account_info(),
                                authority: ctx.accounts.payer.to_account_info(),
                            },
                        ),
                        (listing.price as f64 * ((split.share as u64) as f64 / 100.0)) as u64,
                    )?;
                }
                if split.op_code == 1 {
                    token::transfer(
                        CpiContext::new(
                            ctx.accounts.token_program.to_account_info(),
                            Transfer {
                                from: ctx.accounts.initializer_token_account.to_account_info(),
                                to: ctx.remaining_accounts[i].to_account_info(),
                                authority: ctx.accounts.payer.to_account_info(),
                            },
                        ),
                        (listing.price as f64 * ((split.share as u64) as f64 / 100.0)) as u64,
                    )?;
                }
                i += 1;
            }
        }

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
            ctx.accounts.payer.to_account_info(),
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
            ctx.accounts.payer.to_account_info(),
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
                ctx.accounts.mint.key(),
                *ctx.accounts.payer.key,
                *ctx.accounts.payer.key,
                candy_machine_creator.key(),
                config_line.name, // TODO include mints so name #42069
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
                ctx.accounts.mint.key(),
                candy_machine_creator.key(),
                *ctx.accounts.payer.key,
                *ctx.accounts.metadata.key,
                *ctx.accounts.payer.key,
                Some(candy_machine.data.max_supply),
            ),
            master_edition_infos.as_slice(),
            &[&authority_seeds],
        )?;

        let mut new_update_authority = Some(candy_machine.oracle);

        if !candy_machine.data.retain_authority {
            new_update_authority = Some(ctx.accounts.payer.key());
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

        /*
        let tender_cpx = Transfer {
            from: ctx.accounts.initializer_token_account.to_account_info(),
            to: ctx.accounts.treasury_token_account.to_account_info(),
            authority: ctx.accounts.payer.to_account_info(),
        };
        let tender_cpx_context =
            CpiContext::new(ctx.accounts.token_program.to_account_info(), tender_cpx);
        token::transfer(tender_cpx_context, listing.price)?;
        */

        Ok(())
    }

    pub fn report_batch_cardinalities<'info>(
        ctx: Context<'_, '_, '_, 'info, ReportBatchCardinality<'info>>,
        cardinalities_indices: Vec<Vec<u64>>,
        cardinalities_keys: Vec<String>,
    ) -> Result<(), Error> {
        let batch_cardinalities_report = &mut ctx.accounts.batch_cardinalities_report;
        batch_cardinalities_report.batch_account = ctx.accounts.batch.key();
        batch_cardinalities_report.cardinalities_keys = cardinalities_keys;
        batch_cardinalities_report.cardinalities_indices = cardinalities_indices;

        Ok(())
    }

    pub fn rng_nft_after_quest<'info>(
        ctx: Context<'_, '_, '_, 'info, RngRewardIndiceNFTAfterQuest<'info>>,
        _via_bump: u8,
    ) -> Result<(), Error> {
        const MIN_BATCH_RNG: u8 = 10;
        // the intent is to psuedo-randomly
        // here we need to determine the asset that will be rewarded
        // using the `Reward` structure from a `Quest` account
        // we can certify the provided `reward_mint` account
        // pubkey used by parsing another programs data - ie questing.
        let initializer = &ctx.accounts.initializer;
        let via = &ctx.accounts.via;
        let via_map = &ctx.accounts.via_map;
        let reward_token_account = &ctx.accounts.reward_token_account;
        let batches = &ctx.accounts.batches;
        let quest = &ctx.accounts.quest;
        let questee = &ctx.accounts.questee;
        let reward_ticket = &mut ctx.accounts.reward_ticket;

        if questee.owner != initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if via.token_mint != reward_token_account.mint {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if via.token_mint != via_map.token_mint {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        let reward_indice = quest
            .rewards
            .clone()
            .into_iter()
            .position(|n| n.mint_address == reward_token_account.mint);
        if reward_indice.is_none() {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        // since this is via `ctx.remaining_accounts` - anyone of any intention
        // could bias/specify these accounts in their own bot client.
        // we should rather enforce a minimum batch count and assert that
        // each are unique.
        //
        // when the minimum is unachievable, the number of batches are
        // deferred from the global `batches.counter`.
        //
        // while batches are guaranteed to exceed the minimum rng count
        // required, such that the batch contents are uniformly distributed,
        // and that batch contents are rendered mintable upon upload,
        // the only safety guarantee is a greater sample size to protect
        // against bot inflicted biases as the rng behaviour is still
        // guaranteed random with no predictability.
        let remaining_accounts = 0;
        let ctx_remaining_accounts = ctx.remaining_accounts;
        let mut batch_accounts: Vec<AccountInfo> = Vec::new();
        if batches.counter <= MIN_BATCH_RNG as u64 {
            if ctx_remaining_accounts.len() != batches.counter as usize {
                return Err(QuestError::SuspiciousTransaction.into());
            }
            batch_accounts =
                ctx_remaining_accounts[remaining_accounts..ctx_remaining_accounts.len()].to_vec();
        } else if batches.counter > MIN_BATCH_RNG as u64 {
            if ctx_remaining_accounts.len() != MIN_BATCH_RNG as usize {
                return Err(QuestError::SuspiciousTransaction.into());
            }
            batch_accounts =
                ctx_remaining_accounts[remaining_accounts..MIN_BATCH_RNG as usize].to_vec();
        }

        assert_all_unique_account_infos(&batch_accounts)?;
        let recent_slot_hash = &ctx.accounts.slot_hashes.data.borrow();
        let most_recent = &recent_slot_hash[12..20];
        // nominate for r/shittyprogramming 2022 meme of the year pls
        let rng = u64::from_le_bytes([
            most_recent[0],
            most_recent[1],
            most_recent[2],
            most_recent[3],
            most_recent[4],
            most_recent[5],
            most_recent[6],
            most_recent[7],
        ]);
        let cardinality = via.rarity.clone();
        let (batch_account, cardinality_index) =
            get_valid_batch_for_cardinality(&batch_accounts, rng, cardinality.to_string())?;

        reward_ticket.oracle = batches.oracle;
        reward_ticket.initializer = initializer.key();
        reward_ticket.batch_account = batch_account;
        reward_ticket.cardinality_index = cardinality_index;
        reward_ticket.fulfilled = 0;
        reward_ticket.amount = quest.rewards[reward_indice.unwrap()].amount;
        reward_ticket.reset = false;

        Ok(())
    }

    pub fn recycle_rng_nft_after_quest<'info>(
        ctx: Context<'_, '_, '_, 'info, RecycleRngRewardIndiceNFTAfterQuest<'info>>,
        _via_bump: u8,
    ) -> Result<(), Error> {
        const MIN_BATCH_RNG: u8 = 10;
        // the intent is to psuedo-randomly
        // here we need to determine the asset that will be rewarded
        // using the `Reward` structure from a `Quest` account
        // we can certify the provided `reward_mint` account
        // pubkey used by parsing another programs data - ie questing.
        let initializer = &ctx.accounts.initializer;
        let via = &ctx.accounts.via;
        let via_map = &ctx.accounts.via_map;
        let reward_token_account = &ctx.accounts.reward_token_account;
        let batches = &ctx.accounts.batches;
        let quest = &ctx.accounts.quest;
        let reward_ticket = &mut ctx.accounts.reward_ticket;

        if reward_ticket.amount == 0 {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if via.token_mint != reward_token_account.mint {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if via.token_mint != via_map.token_mint {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if quest
            .rewards
            .clone()
            .into_iter()
            .position(|n| n.mint_address == reward_token_account.mint)
            .is_none()
        {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        // since this is via `ctx.remaining_accounts` - anyone of any intention
        // could bias/specify these accounts in their own bot client.
        // we should rather enforce a minimum batch count and assert that
        // each are unique.
        //
        // when the minimum is unachievable, the number of batches are
        // deferred from the global `batches.counter`.
        //
        // while batches are guaranteed to exceed the minimum rng count
        // required, such that the batch contents are uniformly distributed,
        // and that batch contents are rendered mintable upon upload,
        // the only safety guarantee is a greater sample size to protect
        // against bot inflicted biases as the rng behaviour is still
        // guaranteed random with no predictability.
        let remaining_accounts = 0;
        let ctx_remaining_accounts = ctx.remaining_accounts;
        let mut batch_accounts: Vec<AccountInfo> = Vec::new();
        if batches.counter <= MIN_BATCH_RNG as u64 {
            if ctx_remaining_accounts.len() != batches.counter as usize {
                return Err(QuestError::SuspiciousTransaction.into());
            }
            batch_accounts =
                ctx_remaining_accounts[remaining_accounts..ctx_remaining_accounts.len()].to_vec();
        } else if batches.counter > MIN_BATCH_RNG as u64 {
            if ctx_remaining_accounts.len() != MIN_BATCH_RNG as usize {
                return Err(QuestError::SuspiciousTransaction.into());
            }
            batch_accounts =
                ctx_remaining_accounts[remaining_accounts..MIN_BATCH_RNG as usize].to_vec();
        }

        assert_all_unique_account_infos(&batch_accounts)?;
        let recent_slot_hash = &ctx.accounts.slot_hashes.data.borrow();
        let most_recent = &recent_slot_hash[12..20];
        // nominate for r/shittyprogramming 2022 meme of the year pls
        let rng = u64::from_le_bytes([
            most_recent[0],
            most_recent[1],
            most_recent[2],
            most_recent[3],
            most_recent[4],
            most_recent[5],
            most_recent[6],
            most_recent[7],
        ]);
        let cardinality = via.rarity.clone();
        let (batch_account, cardinality_index) =
            get_valid_batch_for_cardinality(&batch_accounts, rng, cardinality.to_string())?;

        if reward_ticket.oracle == batches.oracle {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if reward_ticket.initializer == initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        reward_ticket.batch_account = batch_account;
        reward_ticket.cardinality_index = cardinality_index;
        reward_ticket.reset = false;

        Ok(())
    }

    pub fn mint_nft_via<'info>(
        ctx: Context<'_, '_, '_, 'info, MintNFTViaRewardTicket<'info>>,
        creator_bump: u8,
        _reward_ticket_bump: u8,
    ) -> Result<(), Error> {
        let now = Clock::get()?.unix_timestamp;

        let via = &mut ctx.accounts.via;
        let mint_hash = &mut ctx.accounts.mint_hash;
        let batch_cardinalities_report = &mut ctx.accounts.batch_cardinalities_report;
        let candy_machine = &mut ctx.accounts.candy_machine;
        let candy_machine_creator = &ctx.accounts.candy_machine_creator;
        let instruction_sysvar_account = &ctx.accounts.instruction_sysvar_account;
        let payer = &mut ctx.accounts.payer;
        let mint = &mut ctx.accounts.mint;

        let reward_ticket = &mut ctx.accounts.reward_ticket;
        if reward_ticket.initializer != payer.key() {
            return Err(QuestError::InvalidInitializer.into());
        }
        if reward_ticket.batch_account != candy_machine.key() {
            return Err(QuestError::SuspiciousCandyMachine.into());
        }
        if reward_ticket.amount == 0 {
            return Err(QuestError::SuspiciousAmounts.into());
        }
        if reward_ticket.reset == true {
            return Err(QuestError::IsReset.into());
        }
        if batch_cardinalities_report.batch_account != candy_machine.key() {
            return Err(QuestError::InvalidCandyMachine.into());
        }

        mint_hash.mint = mint.key();
        mint_hash.minter = payer.key();
        mint_hash.mint_index = via.mints;
        mint_hash.fulfilled = now;

        via.mints += 1;

        // ensures valid via for item in batch
        if via.oracle != ctx.accounts.oracle.key() {
            return Err(QuestError::SuspiciousOracle.into());
        }

        // assert_valid_go_live(payer, clock, candy_machine)?;
        token::mint_to(
            CpiContext::new(
                ctx.accounts.token_program.to_account_info(),
                MintTo {
                    mint: mint.to_account_info(),
                    to: ctx.accounts.mint_ata.to_account_info(),
                    authority: payer.to_account_info(),
                },
            ),
            1,
        )?;

        token::burn(
            CpiContext::new(
                ctx.accounts.token_program.to_account_info(),
                Burn {
                    mint: ctx.accounts.reward_token_mint_account.to_account_info(),
                    from: ctx.accounts.reward_token_account.to_account_info(),
                    authority: payer.to_account_info(),
                },
            ),
            1,
        )?;
        reward_ticket.amount -= 1;

        let recent_slot_hash = &ctx.accounts.slot_hashes.data.borrow();
        let most_recent = &recent_slot_hash[12..20];
        // nominate for r/shittyprogramming 2022 meme of the year pls
        let rng = u64::from_le_bytes([
            most_recent[0],
            most_recent[1],
            most_recent[2],
            most_recent[3],
            most_recent[4],
            most_recent[5],
            most_recent[6],
            most_recent[7],
        ]);
        let config_index: u64 = rng
            % batch_cardinalities_report.cardinalities_indices
                [reward_ticket.cardinality_index as usize]
                .len() as u64;
        let config_line = get_config_line(
            &candy_machine,
            batch_cardinalities_report.cardinalities_indices
                [reward_ticket.cardinality_index as usize][config_index as usize]
                as usize,
        )?;

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
            ctx.accounts.payer.to_account_info(),
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
            ctx.accounts.payer.to_account_info(),
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
                ctx.accounts.mint.key(),
                *ctx.accounts.payer.key,
                *ctx.accounts.payer.key,
                candy_machine_creator.key(),
                config_line.name, // TODO include mints so name #42069
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
                ctx.accounts.mint.key(),
                candy_machine_creator.key(),
                *ctx.accounts.payer.key,
                *ctx.accounts.metadata.key,
                *ctx.accounts.payer.key,
                Some(candy_machine.data.max_supply),
            ),
            master_edition_infos.as_slice(),
            &[&authority_seeds],
        )?;

        let mut new_update_authority = Some(candy_machine.oracle);

        if !candy_machine.data.retain_authority {
            new_update_authority = Some(ctx.accounts.payer.key());
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

        reward_ticket.reset = true;

        Ok(())
    }
}
