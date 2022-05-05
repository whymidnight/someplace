use anchor_lang::prelude::*;
use anchor_lang::Discriminator;
use anchor_spl::token::Mint;
use anchor_spl::token::{self, Burn, Token, TokenAccount, Transfer};
use arrayref::array_ref;
use mpl_token_metadata::state::{
    MAX_CREATOR_LEN, MAX_CREATOR_LIMIT, MAX_NAME_LENGTH, MAX_SYMBOL_LENGTH, MAX_URI_LENGTH,
};
use std::cell::RefMut;

declare_id!("7gAAwQQMSQZf1xA5ZJ9Myq5Zdv2uNmg3KkFEUkAh5ZBy");

#[program]
pub mod questing {
    use super::*;

    pub const QUEST_PIXELBALLZ_SEED: &[u8] = b"quest";
    pub const QUEST_PDA_SEED: &[u8] = b"questing";

    pub fn initialize_quest(ctx: Context<InitQuest>, _quest_index: u64) -> Result<()> {
        let quest_account = &mut ctx.accounts.quest_account;
        quest_account.stage = 0;
        quest_account.start_time = 0;
        quest_account.initializer = ctx.accounts.initializer.key();

        Ok(())
    }
    pub fn burn_balls(ctx: Context<BurnBallz>, _quest_index: u64) -> Result<()> {
        if ctx.accounts.initializer.key() != ctx.accounts.quest_account.initializer {
            return Err(QuestError::InvalidInitializer.into());
        }
        token::burn(
            ctx.accounts.burn_ballz(),
            (2000 as f64 * 10_usize.pow(ctx.accounts.ballz_mint.decimals as u32) as f64) as u64,
        )?;
        let quest_account = &mut ctx.accounts.quest_account;
        quest_account.stage += 1;

        Ok(())
    }
    pub fn transfer_pixelballz(ctx: Context<TransferPixelballz>, _quest_index: u64) -> Result<()> {
        if ctx.accounts.initializer.key() != ctx.accounts.quest_account.initializer {
            return Err(QuestError::InvalidInitializer.into());
        }
        token::transfer(
            ctx.accounts.transfer_pixelballz(),
            (100 as f64 * 10_usize.pow(ctx.accounts.pixelballz_mint.decimals as u32) as f64) as u64,
        )?;
        let quest_account = &mut ctx.accounts.quest_account;
        quest_account.stage += 1;
        quest_account.deposit_token_amount = ctx.accounts.deposit_token_account.key();

        Ok(())
    }

    pub fn start_quest(ctx: Context<StartQuest>, _quest_index: u64) -> Result<()> {
        let now = Clock::get()?.unix_timestamp;
        if ctx.accounts.initializer.key() != ctx.accounts.quest_account.initializer {
            return Err(QuestError::InvalidInitializer.into());
        }
        if ctx.accounts.quest_account.stage != 2 {
            return Err(QuestError::UnexpectedQuestingState.into());
        }
        let quest_account = &mut ctx.accounts.quest_account;
        quest_account.start_time = now;
        quest_account.end_time = now + 2000;
        msg!("{:?} - {:?}", now, now + 2000);

        Ok(())
    }
}

#[derive(Accounts)]
#[instruction(quest_index: u64)]
pub struct StartQuest<'info> {
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub quest_account: Account<'info, QuestAccount>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
pub struct BurnBallz<'info> {
    #[account(mut)]
    pub ballz_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub ballz_mint: Account<'info, Mint>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub quest_account: Account<'info, QuestAccount>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64)]
pub struct TransferPixelballz<'info> {
    #[account(
        init,
        seeds = [&questing::QUEST_PDA_SEED, &questing::QUEST_PIXELBALLZ_SEED, initializer.key().as_ref(), &quest_index.to_le_bytes()],
        bump,
        payer = initializer,
        token::mint = pixelballz_mint,
        token::authority = deposit_token_account
    )]
    pub deposit_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub pixelballz_mint: Account<'info, Mint>,
    #[account(mut)]
    pub pixelballz_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(mut)]
    pub quest_account: Account<'info, QuestAccount>,
    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
#[instruction(quest_index: u64)]
pub struct InitQuest<'info> {
    #[account(mut)]
    pub initializer: Signer<'info>,
    #[account(
        init,
        seeds = [&questing::QUEST_PDA_SEED, initializer.key().as_ref(), &quest_index.to_le_bytes()],
        bump,
        payer = initializer,
        space = QuestAccount::LEN
    )]
    pub quest_account: Account<'info, QuestAccount>,
    pub system_program: Program<'info, System>,
}

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

impl<'info> BurnBallz<'info> {
    fn burn_ballz(&self) -> CpiContext<'_, '_, '_, 'info, Burn<'info>> {
        let cpi_accounts = Burn {
            mint: self.ballz_mint.to_account_info(),
            to: self.ballz_token_account.to_account_info(),
            authority: self.initializer.to_account_info(),
        };
        let cpi_program = self.token_program.to_account_info();
        CpiContext::new(cpi_program, cpi_accounts)
    }
}

impl<'info> TransferPixelballz<'info> {
    fn transfer_pixelballz(&self) -> CpiContext<'_, '_, '_, 'info, Transfer<'info>> {
        let cpi_accounts = Transfer {
            from: self.pixelballz_token_account.to_account_info(),
            to: self.deposit_token_account.to_account_info(),
            authority: self.initializer.to_account_info(),
        };
        let cpi_program = self.token_program.to_account_info();
        CpiContext::new(cpi_program, cpi_accounts)
    }
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
