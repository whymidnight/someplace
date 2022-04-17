use crate::structs::*;
use anchor_lang::prelude::*;

#[account]
pub struct QuestAccount {
    pub stage: u8,
    pub start_time: i64,
    pub end_time: i64,
    pub deposit_token_amount: Pubkey,
    pub initializer: Pubkey,
}

impl QuestAccount {
    pub const LEN: usize = 8 + 8 + 8 + 8 + 32 + 32;
}

#[account]
#[derive(Default)]
pub struct Batch {
    pub name: String, // max 32 bytes
    pub oracle: Pubkey,
    pub data: CandyMachineData,
}

#[account]
#[derive(Default)]
pub struct BatchReceipt {
    pub id: u64,
    pub name: String,
    pub batch_account: Pubkey,
    pub oracle: Pubkey,
}

impl BatchReceipt {
    pub const LEN: usize = 8 + 8 + 32 + 32 + 32;
}

#[account]
#[derive(Default)]
pub struct Batches {
    pub counter: u64,
    pub oracle: Pubkey,
}

impl Batches {
    pub const LEN: usize = 8 + 8 + 32;
}

#[account]
#[derive(Default)]
pub struct TreasuryWhitelist {
    pub whitelist_id: u64,
    pub candy_machine_id: Pubkey,
    pub candy_machine_creator: Pubkey,
    pub treasury_authority: Pubkey,
    pub oracle: Pubkey,
}

impl TreasuryWhitelist {
    pub const LEN: usize = 8 + 8 + 32 + 32 + 32 + 32;
}

#[account]
#[derive(Default)]
pub struct TreasuryAuthority {
    pub whitelists: u64,
    pub oracle: Pubkey,
    pub treasury_decimals: u8,
    pub treasury_token_account: Pubkey,
    pub treasury_mint: Pubkey,
    pub adornment: String, // max 32 bytes
}

impl TreasuryAuthority {
    pub const LEN: usize = 8 + 8 + 32 + 8 + 32 + 32 + 32;
}

#[account]
pub struct Listing {
    pub treasury_authority: Pubkey,
    pub batch: Pubkey,
    pub oracle: Pubkey,
    pub config_index: u64,
    pub price: u64,
    pub lifecycle_start: u64,
    pub lifecycle_end: u64,
}

impl Listing {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 8 + 8 + 8;
}

#[account]
#[derive(Default)]
pub struct Market {
    pub market_decimals: u8,
    pub listings: u64,
    pub name: String, // 32 bytes
    pub market_mint: Pubkey,
    pub market_uid: Pubkey,
    pub oracle: Pubkey,
}

impl Market {
    pub const LEN: usize = 8 + 1 + 8 + 32 + 32 + 32 + 32;
}

#[account]
pub struct MarketListing {
    pub market_authority: Pubkey,
    pub nft_mint: Pubkey,
    pub seller_market_token_account: Pubkey,
    pub index: u64,
    pub price: u64,
    pub fulfilled: u8,
}

impl MarketListing {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 8 + 1;
}
