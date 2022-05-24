use anchor_lang::prelude::*;

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Reward {
    pub mint_address: Pubkey,
    pub rng_threshold: u8,
    // decimals are hardcoded to 0 in ./ix_accounts.rs:AppendQuestReward
    pub amount: u64,
    pub cardinality: Option<String>,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Tender {
    pub mint_address: Pubkey,
    pub amount: u64,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Split {
    pub token_address: Pubkey,
    pub op_code: u8, // 0 - burn, 1 - transfer to `token_address`
    pub share: u8,
}
