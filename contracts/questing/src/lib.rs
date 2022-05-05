use anchor_lang::prelude::*;
use anchor_spl::token::Mint;
use anchor_spl::token::{self, Burn, Token, TokenAccount, Transfer};
use state::*;
use structs::*;

declare_id!("Cv1EGc9jnop3n1fhCNBRSpoQ98uTdGCi4ESAxAwh2Ek5");

mod state;
mod structs;

#[program]
pub mod questing {
    use super::*;

    pub const QUEST_PIXELBALLZ_SEED: &[u8] = b"quest";
    pub const QUEST_PDA_SEED: &[u8] = b"questing";
    pub const QUEST_ORACLE_SEED: &[u8] = b"oracle";
    pub const QUEST_ENTITLEMENT_SEED: &[u8] = b"entitlement";

    pub fn enroll_questor(ctx: Context<EnrollQuestor>) -> Result<()> {
        let questor = &mut ctx.accounts.questor;
        questor.initializer = ctx.accounts.initializer.key();
        questor.quests = 0;

        Ok(())
    }

    pub fn enable_quests(ctx: Context<EnableQuests>) -> Result<()> {
        let quests = &mut ctx.accounts.quests;
        quests.oracle = ctx.accounts.oracle.key();
        quests.quests = 0;

        Ok(())
    }

    pub fn create_quest(
        ctx: Context<CreateQuest>,
        quest_index: u64,
        duration: i64,
        wl_candy_machines: Vec<Pubkey>,
        rewards: Vec<Reward>,
        tender: Option<Tender>,
    ) -> Result<()> {
        let quest = &mut ctx.accounts.quest;
        let quests = &mut ctx.accounts.quests;

        if quests.quests != quest_index {
            return Err(QuestError::UnexpectedQuestingState.into());
        }

        quest.index = quest_index;
        quest.duration = duration;
        quest.oracle = ctx.accounts.oracle.key();
        quest.wl_candy_machines = wl_candy_machines;
        quest.rewards = rewards;
        quest.tender = tender;

        quests.quests += 1;

        Ok(())
    }

    pub fn ammend_quest_with_entitlement(
        ctx: Context<AmmendQuestWithEntitlement>,
        quest_index: u64,
        _quest_bump: u8,
        entitlement: Reward,
    ) -> Result<()> {
        let quest = &mut ctx.accounts.quest;
        if quest.oracle != ctx.accounts.oracle.key() {
            return Err(QuestError::UnexpectedQuestingState.into());
        }
        if quest.index != quest_index {
            return Err(QuestError::UnexpectedQuestingState.into());
        }
        if entitlement.mint_address != ctx.accounts.ballz_mint.key() {
            return Err(QuestError::UnexpectedQuestingState.into());
        }

        quest.entitlement = Some(entitlement);

        Ok(())
    }

    pub fn start_quest(ctx: Context<StartQuest>, _quest_index: u64) -> Result<()> {
        let now = Clock::get()?.unix_timestamp;
        let quest = &mut ctx.accounts.quest;
        let quest_account = &mut ctx.accounts.quest_account;
        quest_account.start_time = now;
        quest_account.end_time = now + 2000;
        msg!("{:?} - {:?}", now, now + 2000);

        token::transfer(
            CpiContext::new(
                ctx.accounts.token_program.to_account_info(),
                Transfer {
                    from: ctx.accounts.pixelballz_token_account.to_account_info(),
                    to: ctx.accounts.deposit_token_account.to_account_info(),
                    authority: ctx.accounts.initializer.to_account_info(),
                },
            ),
            1,
        )?;

        if quest.tender.is_some() {
            let tender = quest.tender.clone().unwrap();
            token::burn(
                CpiContext::new(
                    ctx.accounts.token_program.to_account_info(),
                    Burn {
                        mint: ctx.accounts.ballz_mint.to_account_info(),
                        to: ctx.accounts.ballz_token_account.to_account_info(),
                        authority: ctx.accounts.initializer.to_account_info(),
                    },
                ),
                tender.amount,
            )?;
        }

        Ok(())
    }

    pub fn end_quest(
        ctx: Context<EndQuest>,
        quest_index: u64,
        deposit_token_account_bump: u8,
    ) -> Result<()> {
        let now = Clock::get()?.unix_timestamp;
        let quest_account = &mut ctx.accounts.quest_account;
        let initializer = &ctx.accounts.initializer.key();
        quest_account.start_time = now;
        quest_account.end_time = now + 2000;
        msg!("{:?} - {:?}", now, now + 2000);
        let seeds = &[
            QUEST_PDA_SEED,
            QUEST_PIXELBALLZ_SEED,
            initializer.as_ref(),
            &quest_index.to_le_bytes(),
        ];
        let (deposit_token_account, bump) = Pubkey::find_program_address(seeds, &crate::ID);
        if deposit_token_account != ctx.accounts.deposit_token_account.key() {
            return Err(QuestError::InvalidInitializer.into());
        }
        if deposit_token_account_bump != bump {
            return Err(QuestError::InvalidInitializer.into());
        }
        let seeds_with_bump = &[
            QUEST_PDA_SEED,
            QUEST_PIXELBALLZ_SEED,
            initializer.as_ref(),
            &quest_index.to_le_bytes(),
            &bump.to_le_bytes(),
        ];
        let deposit_token_account_authority = &[&seeds_with_bump[..]];

        token::transfer(
            CpiContext::new_with_signer(
                ctx.accounts.token_program.to_account_info(),
                Transfer {
                    from: ctx.accounts.deposit_token_account.to_account_info(),
                    to: ctx.accounts.pixelballz_token_account.to_account_info(),
                    authority: ctx.accounts.deposit_token_account.to_account_info(),
                },
                deposit_token_account_authority,
            ),
            1,
        )?;

        Ok(())
    }
}

#[derive(Accounts)]
#[instruction(quest_index: u64)]
pub struct StartQuest<'info> {
    #[account(mut)]
    pub quest: Box<Account<'info, Quest>>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub ballz_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub ballz_mint: Box<Account<'info, Mint>>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), QUEST_PIXELBALLZ_SEED.as_ref(), initializer.key().as_ref(), quest.key().as_ref()],
        bump,
        payer = initializer,
        token::mint = pixelballz_mint,
        token::authority = deposit_token_account
    )]
    pub deposit_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    #[account(mut)]
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    #[account(
        init,
        seeds = [QUEST_PDA_SEED.as_ref(), initializer.key().as_ref(), quest.key().as_ref()],
        bump,
        payer = initializer,
        space = QuestAccount::LEN
    )]
    pub quest_account: Box<Account<'info, QuestAccount>>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
pub struct EndQuest<'info> {
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub deposit_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub pixelballz_mint: Box<Account<'info, Mint>>,
    #[account(mut)]
    pub pixelballz_token_account: Box<Account<'info, TokenAccount>>,
    #[account(mut)]
    pub quest_account: Box<Account<'info, QuestAccount>>,
    pub token_program: Program<'info, Token>,
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

#[error]
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
}
