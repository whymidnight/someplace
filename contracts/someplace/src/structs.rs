use anchor_lang::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct Creator {
    pub address: Pubkey,
    pub verified: bool,
    pub share: u8,
}
/// Configurations options for the gatekeeper.
#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct GatekeeperConfig {
    pub gatekeeper_network: Pubkey,
    pub expire_on_use: bool,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub enum EndSettingType {
    Date,
    Amount,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct EndSettings {
    pub end_setting_type: EndSettingType,
    pub number: u64,
}
#[derive(AnchorSerialize, AnchorDeserialize, Clone, Default)]
pub struct HiddenSettings {
    pub name: String,
    pub uri: String,
    pub hash: [u8; 32],
}
#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct WhitelistMintSettings {
    pub mode: WhitelistMintMode,
    pub mint: Pubkey,
    pub presale: bool,
    pub discount_price: Option<u64>,
}

#[derive(AnchorSerialize, AnchorDeserialize, Clone, PartialEq)]
pub enum WhitelistMintMode {
    BurnEveryTime,
    NeverBurn,
    DeferCreator,
}

/// Candy machine settings data.
#[derive(AnchorSerialize, AnchorDeserialize, Clone, Default)]
pub struct CandyMachineData {
    pub uuid: String,
    pub price: u64,
    pub symbol: String,
    pub seller_fee_basis_points: u16,
    pub max_supply: u64,
    pub is_mutable: bool,
    pub retain_authority: bool,
    pub go_live_date: Option<i64>,
    pub end_settings: Option<EndSettings>,
    pub creators: Vec<Creator>,
    pub hidden_settings: Option<HiddenSettings>,
    pub whitelist_mint_settings: Option<WhitelistMintSettings>,
    pub items_available: u64,
    pub gatekeeper: Option<GatekeeperConfig>,
}

#[derive(AnchorSerialize, AnchorDeserialize, Serialize, Deserialize, Debug)]
pub struct ConfigLine {
    pub name: String,
    pub cardinality: String,
    pub uri: String,
}

#[derive(AnchorSerialize, AnchorDeserialize, Serialize, Deserialize, Clone)]
pub struct Split {
    pub token_address: Pubkey,
    pub op_code: u8, // 0 - burn, 1 - transfer to `token_address`
    pub share: u8,
}
