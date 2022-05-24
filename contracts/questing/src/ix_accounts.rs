use crate::constants::*;
use crate::state::*;
use anchor_lang::prelude::*;
use anchor_spl::associated_token::AssociatedToken;
use anchor_spl::token::{Mint, Token, TokenAccount};

#[derive(Accounts)]
pub struct StartQuest<'info> {
    #[account(mut)]
    pub quest: Box<Account<'info, Quest>>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), questee.key().as_ref(), quest.key().as_ref()],
        bump,
        payer = initializer,
        token::mint = pixelballz_mint,
        token::authority = deposit_token_account
    )]
    pub deposit_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    #[account(mut)]
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), questor.key().as_ref(), questee.key().as_ref(), quest.key().as_ref()],
        bump,
        payer = initializer,
        space = QuestAccount::LEN
    )]
    pub quest_acc: Account<'info, QuestAccount>,
    #[account(mut)]
    pub questor: Account<'info, Questor>,
    #[account(mut)]
    pub questee: Account<'info, Questee>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
#[instruction(deposit_token_account_bump: u8)]
pub struct EndQuest<'info> {
    #[account(
        init,
        seeds = [QUEST_REWARD_SEED.as_ref(), questor.key().as_ref(), questee.key().as_ref(), quest.key().as_ref()],
        bump,
        payer = initializer,
        space = QuestQuesteeEndReceipt::LEN
    )]
    pub quest_questee_receipt: Account<'info, QuestQuesteeEndReceipt>,
    #[account(mut)]
    pub quest: Account<'info, Quest>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(
        mut,
        seeds = [QUEST_PDA_SEED.as_ref(), questee.key().as_ref(), quest.key().as_ref()],
        bump=deposit_token_account_bump,
    )]
    pub deposit_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    #[account(mut)]
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub quest_acc: Box<Account<'info, QuestAccount>>,
    #[account(mut)]
    pub questor: Box<Account<'info, Questor>>,
    #[account(mut)]
    pub questee: Account<'info, Questee>,
    /// CHECK: am lazy
    pub oracle: UncheckedAccount<'info>,
    pub token_program: Program<'info, Token>,
    pub associated_token_program: Program<'info, AssociatedToken>,
    /// CHECK: am lazy
    pub slot_hashes: UncheckedAccount<'info>,
    pub system_program: Program<'info, System>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64)]
pub struct CreateQuest<'info> {
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(
        init,
        seeds = [QUEST_ORACLE_SEED.as_ref(), oracle.key().as_ref(), &quest_index.to_le_bytes()],
        bump,
        payer = oracle,
        space = Quest::LEN
    )]
    pub quest: Box<Account<'info, Quest>>,
    #[account(mut)]
    pub quests: Box<Account<'info, Quests>>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64, quest_bump: u8)]
pub struct AppendQuestReward<'info> {
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(mut, seeds = [QUEST_ORACLE_SEED.as_ref(), oracle.key().as_ref(), &quest_index.to_le_bytes()], bump=quest_bump)]
    pub quest: Box<Account<'info, Quest>>,
    #[account(
        init,
        payer = oracle,
        mint::decimals = 0,
        mint::authority = quest,
    )]
    pub reward_mint: Account<'info, Mint>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64, quest_bump: u8)]
pub struct ResetQuestRewards<'info> {
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(mut, seeds = [QUEST_ORACLE_SEED.as_ref(), oracle.key().as_ref(), &quest_index.to_le_bytes()], bump=quest_bump)]
    pub quest: Box<Account<'info, Quest>>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64, quest_bump: u8)]
pub struct AmmendQuestWithEntitlement<'info> {
    #[account(mut)]
    pub ballz_mint: Box<Account<'info, Mint>>,
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(mut, seeds = [QUEST_ORACLE_SEED.as_ref(), oracle.key().as_ref(), &quest_index.to_le_bytes()], bump=quest_bump)]
    pub quest: Box<Account<'info, Quest>>,
    #[account(
        init,
        seeds = [QUEST_ORACLE_SEED.as_ref(), QUEST_ENTITLEMENT_SEED.as_ref(), oracle.key().as_ref(), &quest_index.to_le_bytes()],
        bump,
        payer = oracle,
        token::mint = ballz_mint,
        token::authority = entitlement_token_account
    )]
    pub entitlement_token_account: Box<Account<'info, TokenAccount>>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
pub struct EnableQuests<'info> {
    #[account(mut)]
    pub oracle: Signer<'info>,
    #[account(
        init,
        seeds = [QUEST_ORACLE_SEED.as_ref(), oracle.key().as_ref()],
        bump,
        payer = oracle,
        space = Quests::LEN
    )]
    pub quests: Box<Account<'info, Quests>>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct EnrollQuestor<'info> {
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), initializer.key().as_ref()],
        bump,
        payer = initializer,
        space = Questor::LEN
    )]
    pub questor: Box<Account<'info, Questor>>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct EnrollQuestee<'info> {
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), pixelballz_mint.key().as_ref()],
        bump,
        payer = initializer,
        space = Questee::LEN
    )]
    pub questee: Box<Account<'info, Questee>>,
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    pub questor: Box<Account<'info, Questor>>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct UpdateQuestee<'info> {
    #[account(mut)]
    pub new_owner: Signer<'info>,
    #[account(
        mut,
        has_one = owner,
    )]
    pub questee: Box<Account<'info, Questee>>,
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    /// CHECK: am lazy
    pub owner: UncheckedAccount<'info>,
    #[account(mut)]
    pub questor: Box<Account<'info, Questor>>,
    pub system_program: Program<'info, System>,
}
