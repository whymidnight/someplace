use anchor_lang::prelude::*;

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Reward {
    pub mint_address: Pubkey,
    pub rng_threshold: Option<u8>,
    pub amount: u8,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Tender {
    pub mint_address: Pubkey,
    pub amount: u64,
}
