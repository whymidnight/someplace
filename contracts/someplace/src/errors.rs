use anchor_lang::prelude::*;

#[error_code]
pub enum QuestError {
    #[msg("Unexpected questing state")]
    UnexpectedQuestingState,
    #[msg("Invalid initizalizer")]
    InvalidInitializer,
    #[msg("Is timelocked")]
    IsTimelocked,
    #[msg("Numerical overflow error!")]
    NumericalOverflowError,
    #[msg("Index greater than length!")]
    IndexGreaterThanLength,
    #[msg("Unable to find an unused config line near your random number index")]
    CannotFindUsableConfigLine,
    #[msg("Uuid must be exactly of 6 length")]
    UuidMustBeExactly6Length,
    #[msg("Invalid string")]
    InvalidString,
    #[msg("Invalid account data")]
    InvalidAccountData,
    #[msg("Suspicious Accounts")]
    SuspiciousAccounts,
    #[msg("Suspicious Transaction")]
    SuspiciousTransaction,
    #[msg("Suspicious Treasury")]
    SuspiciousTreasury,
    #[msg("Suspicious Treasury Mint")]
    SuspiciousTreasuryMint,
    #[msg("Suspicious Candy Machine")]
    SuspiciousCandyMachine,
    #[msg("Suspicious Amounts")]
    SuspiciousAmounts,
    #[msg("Is Reset")]
    IsReset,
    #[msg("Invalid Candy Batch")]
    InvalidCandyMachine,
    #[msg("Suspicious Oracle")]
    SuspiciousOracle,
    #[msg("Suspicious Token Mint")]
    SuspiciousTokenMint,
    #[msg("Suspicious Via Token Mint")]
    SuspiciousViaTokenMint,
    #[msg("Malformed Reward Mint")]
    MalformedRewardMint,
    #[msg("Suspicious Batches Length")]
    SuspiciousBatchesLength,
    #[msg("Invalid Amount")]
    InvalidAmount,
}
