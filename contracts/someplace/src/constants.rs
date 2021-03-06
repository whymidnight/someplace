use mpl_token_metadata::state::{
    MAX_CREATOR_LEN, MAX_CREATOR_LIMIT, MAX_NAME_LENGTH, MAX_SYMBOL_LENGTH, MAX_URI_LENGTH,
};

// const DURATION: i64 = 60 * 5;
// NAME[32] +
pub const MAX_CARDINALITY_LENGTH: usize = 32;
pub const CONFIG_LINE_SIZE: usize =
    4 + MAX_NAME_LENGTH + 4 + MAX_CARDINALITY_LENGTH + 4 + MAX_URI_LENGTH;
pub const CONFIG_ARRAY_START: usize = 8 + // key
    32 + // authority
    32 + //wallet
    33 + // token mint
    4 + 6 + // uuid
    8 + // price
    8 + // items available
    9 + // go live
    10 + // end settings
    4 + MAX_SYMBOL_LENGTH + // u32 len + symbol
    2 + // seller fee basis points
    4 + MAX_CREATOR_LIMIT*MAX_CREATOR_LEN + // optional + u32 len + actual vec
    8 + //max supply
    1 + // is mutable
    1 + // retain authority
    1 + // option for hidden setting
    4 + MAX_NAME_LENGTH + // name length,
    4 + MAX_URI_LENGTH + // uri length,
    32 + // hash
    4 +  // max number of lines;
    8 + // items redeemed
    1 + // whitelist option
    1 + // whitelist mint mode
    1 + // allow presale
    9 + // discount price
    32 + // mint key for whitelist
    1 + 32 + 1 // gatekeeper
;
pub const PREFIX: &[u8] = b"someplace";
pub const LISTING: &[u8] = b"publiclisting";
pub const LISTINGTOKEN: &[u8] = b"listingtoken";
pub const MARKET: &[u8] = b"market";
pub const BENEFIT_TOKEN: &[u8] = b"ballz";
pub const TREASURY_MINT: &[u8] = b"treasury_mint";
pub const TREASURY_WHITELIST: &[u8] = b"treasury_whitelist";
pub const MINTYHASH: &[u8] = b"mintyhash";
pub const VIA: &[u8] = b"via";
pub const BATCH_CARDINALITIES: &[u8] = b"batch_cardinalities";
pub const VIA_MINT_HASH: &[u8] = b"via_mint_hash";
