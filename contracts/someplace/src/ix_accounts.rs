use crate::constants::*;
use crate::helper_fns::*;
use crate::state::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use anchor_lang::solana_program::sysvar;
use anchor_spl::token::{Mint, Token, TokenAccount};

#[derive(Accounts)]
#[instruction(config_index: u64)]
pub struct CreateListing<'info> {
    #[account(mut, has_one = oracle)]
    pub batch: Account<'info, Batch>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(
        init,
        seeds = [oracle.key().as_ref(), batch.key().as_ref(), config_index.to_le_bytes().as_ref()],
        bump,
        payer = oracle,
        space = Listing::LEN
    )]
    pub listing: Account<'info, Listing>,
    pub treasury_authority: Account<'info, TreasuryAuthority>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(config_index: u64)]
pub struct ModifyListing<'info> {
    #[account(mut, has_one = oracle)]
    pub batch: Account<'info, Batch>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(
        seeds = [oracle.key().as_ref(), batch.key().as_ref(), config_index.to_le_bytes().as_ref()],
        bump,
    )]
    pub listing: Account<'info, Listing>,
    pub treasury_authority: Account<'info, TreasuryAuthority>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(candy_machine_creator: Pubkey)]
pub struct AddWhitelistedCM<'info> {
    #[account(
        init,
        seeds = [oracle.key().as_ref(), TREASURY_WHITELIST.as_ref(), treasury_authority.key().as_ref(), candy_machine_creator.as_ref()],
        bump,
        payer = oracle,
        space = TreasuryWhitelist::LEN
    )]
    pub treasury_whitelist: Account<'info, TreasuryWhitelist>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(mut)]
    pub treasury_authority: Account<'info, TreasuryAuthority>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(market_uid: Pubkey)]
pub struct InitMarket<'info> {
    #[account(
        init,
        seeds = [PREFIX.as_ref(), MARKET.as_ref(), oracle.key().as_ref(), market_uid.as_ref()],
        bump,
        payer = oracle,
        space = Market::LEN,
    )]
    pub market_authority: Account<'info, Market>,
    #[account(mut)]
    pub market_mint: Account<'info, Mint>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(index: u64)]
pub struct InitMarketListing<'info> {
    #[account(mut)]
    pub market_authority: Account<'info, Market>,
    #[account(
        init,
        seeds = [PREFIX.as_ref(), LISTING.as_ref(), market_authority.key().as_ref(), market_authority.listings.to_le_bytes().as_ref()],
        bump,
        payer = seller,
        space = MarketListing::LEN,
    )]
    pub market_listing: Account<'info, MarketListing>,
    #[account(
        init,
        seeds = [PREFIX.as_ref(), LISTINGTOKEN.as_ref(), market_authority.key().as_ref(), index.to_le_bytes().as_ref()],
        bump,
        payer = seller,
        token::mint = nft_mint,
        token::authority = market_authority
    )]
    pub market_listing_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub seller: Signer<'info>,
    #[account(mut)]
    pub seller_nft_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub seller_market_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub nft_mint: Account<'info, Mint>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}
#[derive(Accounts)]
pub struct FulfillMarketListing<'info> {
    #[account(mut)]
    pub market_authority: Account<'info, Market>,
    #[account(mut)]
    pub market_listing: Account<'info, MarketListing>,
    #[account(mut)]
    pub market_listing_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub buyer: Signer<'info>,
    #[account(mut)]
    pub nft_mint: Account<'info, Mint>,
    #[account(mut)]
    pub buyer_nft_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub buyer_market_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub seller_market_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    /// CHECK: !islazy && /s
    pub oracle: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
}
#[derive(Accounts)]
pub struct UnlistMarketListing<'info> {
    #[account(mut)]
    pub market_authority: Account<'info, Market>,
    #[account(mut)]
    pub market_listing: Account<'info, MarketListing>,
    #[account(mut)]
    pub market_listing_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub seller: Signer<'info>,
    #[account(mut)]
    pub nft_mint: Account<'info, Mint>,
    #[account(mut)]
    pub seller_nft_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    /// CHECK: !islazy && /s
    pub oracle: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
pub struct InitTreasury<'info> {
    #[account(
        init,
        seeds = [PREFIX.as_ref(), BENEFIT_TOKEN.as_ref(), oracle.key().as_ref()],
        bump,
        payer = oracle,
        space = TreasuryAuthority::LEN
    )]
    pub treasury_authority: Account<'info, TreasuryAuthority>,
    #[account(
        init,
        seeds = [PREFIX.as_ref(), BENEFIT_TOKEN.as_ref(), TREASURY_MINT.as_ref(), oracle.key().as_ref()],
        bump,
        payer = oracle,
        token::mint = treasury_token_mint,
        token::authority = treasury_authority
    )]
    pub treasury_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub treasury_token_mint: Account<'info, Mint>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(mut)]
    pub oracle_token_account: Account<'info, TokenAccount>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
pub struct SellFor<'info> {
    #[account(mut)]
    pub depo_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub depo_mint: Account<'info, Mint>,
    #[account(mut)]
    /// CHECK: legacy
    pub metadata: UncheckedAccount<'info>,
    #[account(mut)]
    pub treasury_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub treasury_token_mint: Account<'info, Mint>,
    #[account(mut)]
    pub treasury_authority: Account<'info, TreasuryAuthority>,
    #[account(mut)]
    pub treasury_whitelist: Account<'info, TreasuryWhitelist>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub initializer_token_account: Box<Account<'info, TokenAccount>>,
    /// CHECK: legacy
    pub oracle: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
pub struct EnableBatches<'info> {
    #[account(
        init,
        seeds = [oracle.key().as_ref()],
        bump,
        payer = oracle,
        space = Batches::LEN
    )]
    pub batches: Account<'info, Batches>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(data: CandyMachineData)]
pub struct NewBatch<'info> {
    /// CHECK: legacy
    #[account(zero, rent_exempt = skip, constraint = batch_account.to_account_info().owner == program_id && batch_account.to_account_info().data_len() >= get_space_for_batch(data)?)]
    pub batch_account: UncheckedAccount<'info>,
    #[account(mut)]
    pub batches: Account<'info, Batches>,
    #[account(
        init,
        seeds = [oracle.key().as_ref(), batches.counter.to_le_bytes().as_ref()],
        bump,
        payer = oracle,
        space = BatchReceipt::LEN
    )]
    pub batch_receipt: Account<'info, BatchReceipt>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    pub system_program: Program<'info, System>,
}

/// Add multiple config lines to the candy machine.
#[derive(Accounts)]
pub struct Sync<'info> {
    #[account(mut, has_one = oracle)]
    pub batch: Account<'info, Batch>,
    pub oracle: Signer<'info>,
}

/// Mint a new NFT pseudo-randomly from the config array.
#[derive(Accounts)]
#[instruction(creator_bump: u8)]
pub struct MintNFT<'info> {
    #[account(mut)]
    pub listing: Box<Account<'info, Listing>>,
    #[account(
        init,
        seeds = [oracle.key().as_ref(), listing.key().as_ref(), listing.mints.to_le_bytes().as_ref()],
        bump,
        payer = payer,
        space = MintHash::LEN
    )]
    pub mint_hash: Box<Account<'info, MintHash>>,
    #[account(
    mut,
    has_one = oracle
    )]
    pub candy_machine: Box<Account<'info, Batch>>,
    #[account(seeds=[PREFIX.as_ref(), candy_machine.key().as_ref()], bump=creator_bump)]
    /// CHECK: legacy
    pub candy_machine_creator: UncheckedAccount<'info>,
    #[account(mut)]
    pub payer: Signer<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub oracle: UncheckedAccount<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub metadata: UncheckedAccount<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub mint: UncheckedAccount<'info>,
    pub mint_authority: Signer<'info>,
    pub update_authority: Signer<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub master_edition: UncheckedAccount<'info>,
    #[account(address = mpl_token_metadata::id())]
    /// CHECK: legacy
    pub token_metadata_program: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
    pub rent: Sysvar<'info, Rent>,
    pub clock: Sysvar<'info, Clock>,
    #[account(address = sysvar::instructions::id())]
    /// CHECK: legacy
    pub instruction_sysvar_account: UncheckedAccount<'info>,
    #[account(mut)]
    pub treasury_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub initializer_token_account: Box<Account<'info, TokenAccount>>,
}

/// Mint a new NFT pseudo-randomly from the config array.
#[derive(Accounts)]
#[instruction(creator_bump: u8)]
pub struct MintNFTRarity<'info> {
    #[account(
    mut,
    has_one = oracle
    )]
    pub candy_machine: Account<'info, Batch>,
    #[account(seeds=[PREFIX.as_ref(), candy_machine.key().as_ref()], bump=creator_bump)]
    /// CHECK: legacy
    pub candy_machine_creator: UncheckedAccount<'info>,
    pub payer: Signer<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub oracle: UncheckedAccount<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub metadata: UncheckedAccount<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub mint: UncheckedAccount<'info>,
    #[account(mut)]
    /// CHECK: legacy
    pub master_edition: UncheckedAccount<'info>,
    #[account(address = mpl_token_metadata::id())]
    /// CHECK: legacy
    pub token_metadata_program: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
    pub rent: Sysvar<'info, Rent>,
    pub clock: Sysvar<'info, Clock>,
    #[account(address = sysvar::instructions::id())]
    /// CHECK: legacy
    pub instruction_sysvar_account: UncheckedAccount<'info>,
    #[account(mut)]
    pub treasury_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub initializer_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    /// CHECK: legacy
    pub nft_token_account: UncheckedAccount<'info>,
    /// CHECK: legacy
    pub recent_blockhashes: UncheckedAccount<'info>,
}
