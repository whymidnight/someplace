use crate::structs::*;
use anchor_lang::prelude::*;
use borsh::BorshDeserialize;
use std::result::Result;

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

impl Batch {
    pub fn from_account_info(a: &AccountInfo) -> Result<Batch, Error> {
        Ok(BorshDeserialize::deserialize(
            &mut a.data.borrow_mut().as_ref(),
        )?)
    }
}

#[account]
#[derive(Default)]
pub struct BatchCardinalitiesReport {
    pub batch_account: Pubkey,
    pub cardinalities_indices: Vec<Vec<u64>>,
    pub cardinalities_keys: Vec<String>,
}

impl BatchCardinalitiesReport {
    pub fn from_account_info(a: &AccountInfo) -> Result<BatchCardinalitiesReport, Error> {
        Ok(BorshDeserialize::deserialize(
            &mut a.data.borrow()[8..a.data_len()].as_ref(),
        )?)
    }
    pub fn get_space(cardinalities: Vec<Vec<u64>>) -> usize {
        let mut sum_indices_mem_size: usize = 0;
        let mut sum_keys_mem_size: usize = 0;
        sum_indices_mem_size += 4 * cardinalities.len();
        sum_keys_mem_size += 4 * cardinalities.len();
        for cardinality in cardinalities {
            sum_indices_mem_size += 8 * cardinality.len();
            sum_keys_mem_size += 32 * cardinality.len();
        }

        8 + 4 + sum_indices_mem_size + 4 + sum_keys_mem_size + 32
    }
}

#[account]
#[derive(Default)]
pub struct BatchReceipt {
    pub id: u64,
    pub name: String,
    pub batch_account: Pubkey,
    pub oracle: Pubkey,
    pub items: u64,
}

impl BatchReceipt {
    pub const LEN: usize = 8 + 8 + 32 + 32 + 32 + 8;
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
    pub splits: Vec<Split>,
    // via_mints vector represents a list
    // of whitelisted token mints from
    // disparate systems to qualify
    // rarer/specific types of rarity
    // assets. this is only permissable
    // by the `oracle` pubkey to guard and
    // permit aforementioned whitelisted mints
    // for such rarity.
    //
    //
    // instead of imposing this as a unique pda,
    // we store this in here to reduce the
    // amount of accounts specified in a
    // mint_nft instruction since it effectively
    // subtracts the account needed for this data since
    // the data is already embedded in an included
    // account.
    pub via_mints: Vec<ViaMint>,
    pub adornment: String,
}

impl TreasuryAuthority {
    pub const LEN: usize = 8
        + 8
        + 32
        + 8
        + 32
        + 32
        + (4 + (std::mem::size_of::<Split>() * 10)) // max of 10 splits
        + (4 + (std::mem::size_of::<ViaMint>() * 10)) // max of 10 via(wl) mints
        + 32 // max 32 bytes
    ;
}

#[account]
pub struct Listing {
    pub treasury_authority: Pubkey,
    pub batch: Pubkey,
    pub oracle: Pubkey,
    pub config_index: u64,
    pub price: u64,
    pub lifecycle_start: u64,
    pub is_listed: bool,
    pub mints: u64,
}

impl Listing {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 8 + 8 + 2 + 8;
}

#[account]
pub struct Via {
    pub oracle: Pubkey,
    pub treasury_authority: Pubkey,
    pub token_mint: Pubkey,
    pub mints: u64,
    pub rarity: String, // 32 byte max
}

impl Via {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 32;
}

#[account]
pub struct ViaMapping {
    pub token_mint: Pubkey,
    pub vias_index: u64,
}

impl ViaMapping {
    pub const LEN: usize = 8 + 32 + 8;
}

#[account]
pub struct Vias {
    pub oracle: Pubkey,
    pub treasury_authority: Pubkey,
    pub vias: u64,
}

impl Vias {
    pub const LEN: usize = 8 + 32 + 32 + 8;
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
    pub seller: Pubkey,
    pub seller_market_token_account: Pubkey,
    pub index: u64,
    pub price: u64,
    pub listed_at: u64,
    pub fulfilled: i64,
}

impl MarketListing {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 32 + 8 + 8 + 8 + 8;
}

#[account]
pub struct MintHash {
    pub mint: Pubkey,
    pub minter: Pubkey,
    pub batch: Pubkey,
    pub config_index: u64,
    pub mint_index: u64,
    pub fulfilled: i64,
}

impl MintHash {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 8 + 8;
}

#[account]
pub struct RewardTicket {
    pub oracle: Pubkey,
    pub initializer: Pubkey,
    pub batch_account: Pubkey,
    // pub config_index: u64,
    // `config_index` was considered but left out to
    // obscure the reward a little from scraping but
    // also to optimize +1 batch accounts in
    // a `rng_nft_after_quest` instruction. since
    // psuedo-random `config_index` can be derived
    // later in a `mint_nft_via_reward_ticket` ix.
    pub fulfilled: i64,
    pub cardinality_index: u64,
    pub amount: u64,
    pub reset: bool,
}

impl RewardTicket {
    pub const LEN: usize = 8 + 32 + 32 + 32 + 8 + 8 + 8 + 2;
}
