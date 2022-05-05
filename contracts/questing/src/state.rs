use crate::structs::*;
use anchor_lang::prelude::*;

#[account]
pub struct QuestAccount {
    pub start_time: i64,
    pub end_time: i64,
    pub deposit_token_amount: Pubkey,
    pub initializer: Pubkey,
}

impl QuestAccount {
    pub const LEN: usize = 8 + 8 + 8 + 32 + 32;
}

#[account]
pub struct Quests {
    pub oracle: Pubkey,
    pub quests: u64,
}

impl Quests {
    pub const LEN: usize = 8 + 8 + 32;
}

#[account]
pub struct Quest {
    pub index: u64,
    pub duration: i64,
    pub oracle: Pubkey,
    pub wl_candy_machines: Vec<Pubkey>,
    pub entitlement: Option<Reward>,
    pub rewards: Vec<Reward>,
    pub tender: Option<Tender>,
}

impl Quest {
    pub const LEN: usize = 8
        + 8
        + 8
        + 32
        + (4 + (10 * 32))
        + std::mem::size_of::<Reward>()
        + (4 + (10 * std::mem::size_of::<Reward>()))
        + (std::mem::size_of::<Tender>());
}

#[account]
pub struct Questor {
    pub initializer: Pubkey,
    pub quests: u64,
}

impl Questor {
    pub const LEN: usize = 8 + 32 + 8;
}
