use anchor_lang::prelude::*;
use anchor_spl::associated_token::{self, Create};
use anchor_spl::token::{self, Burn, MintTo, Transfer};
use constants::*;
use errors::QuestError;
use helper_fns::*;
use ix_accounts::*;
use std::result::Result;
use structs::*;

declare_id!("Cr4keTx8UQiQ5F9TzTGdQ5dkcMHjxhYSAaHkHhUSABCk");

mod constants;
mod errors;
pub mod helper_fns;
mod ix_accounts;
pub mod state;
pub mod structs;

#[program]
pub mod questing {
    use super::*;

    pub fn enroll_questor(ctx: Context<EnrollQuestor>) -> Result<(), Error> {
        let questor = &mut ctx.accounts.questor;
        questor.initializer = ctx.accounts.initializer.key();
        questor.xp = 0;

        Ok(())
    }

    pub fn enroll_questee(ctx: Context<EnrollQuestee>) -> Result<(), Error> {
        let questee = &mut ctx.accounts.questee;
        let initializer = &ctx.accounts.initializer;
        let questor = &ctx.accounts.questor;
        let pixelballz_mint = &ctx.accounts.pixelballz_mint;
        let pixelballz_token_account = &ctx.accounts.pixelballz_token_account;

        if pixelballz_token_account.mint != pixelballz_mint.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if pixelballz_token_account.owner != initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if pixelballz_token_account.amount != 1 {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if questor.initializer != initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        questee.owner = questor.initializer;
        questee.pixelballz_mint = ctx.accounts.pixelballz_mint.key();
        questee.quests = 0;
        questee.xp = 0;

        Ok(())
    }

    pub fn update_questee(ctx: Context<UpdateQuestee>) -> Result<(), Error> {
        let questee = &mut ctx.accounts.questee;
        let initializer = &ctx.accounts.new_owner;
        let pixelballz_mint = &ctx.accounts.pixelballz_mint;
        let pixelballz_token_account = &ctx.accounts.pixelballz_token_account;

        if pixelballz_token_account.mint != pixelballz_mint.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if pixelballz_token_account.owner != initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if pixelballz_token_account.amount != 1 {
            return Err(QuestError::SuspiciousTransaction.into());
        }

        questee.owner = initializer.key();

        Ok(())
    }

    pub fn enable_quests(ctx: Context<EnableQuests>) -> Result<(), Error> {
        let quests = &mut ctx.accounts.quests;
        quests.oracle = ctx.accounts.oracle.key();
        quests.quests = 0;

        Ok(())
    }

    pub fn create_quest(
        ctx: Context<CreateQuest>,
        quest_index: u64,
        name: String,
        duration: i64,
        wl_candy_machines: Vec<Pubkey>,
        tender: Option<Tender>,
        tender_splits: Option<Vec<Split>>,
        xp: u64,
        required_level: Option<u64>,
        enabled: bool,
    ) -> Result<(), Error> {
        let quest = &mut ctx.accounts.quest;
        let quests = &mut ctx.accounts.quests;

        if quests.quests != quest_index {
            return Err(QuestError::UnexpectedQuestingState.into());
        }

        quest.enabled = enabled;
        quest.index = quest_index;
        quest.name = name;
        quest.duration = duration;
        quest.oracle = ctx.accounts.oracle.key();
        quest.wl_candy_machines = wl_candy_machines;
        quest.rewards = Vec::new();
        if tender.is_some() && tender_splits.is_some() {
            quest.tender = tender;
            quest.tender_splits = tender_splits;
        }
        quest.xp = xp;

        if required_level.is_some() {
            quest.required_level = required_level.unwrap();
            quest.required_xp = get_sum_xp_from_level(quest.required_level).unwrap();
        } else {
            quest.required_level = 0;
            quest.required_xp = 0;
        }

        quests.quests += 1;

        Ok(())
    }

    pub fn register_quest_reward(
        ctx: Context<AppendQuestReward>,
        _quest_index: u64,
        _quest_bump: u8,
        reward: Reward,
    ) -> Result<(), Error> {
        let mut reward_params = reward;
        let quest = &mut ctx.accounts.quest;

        // assert valid reward threshold sum before appending to rewards
        assert_valid_sum_rng_thresholds(quest.rewards.clone(), reward_params.rng_threshold)
            .unwrap();

        reward_params = Reward {
            mint_address: ctx.accounts.reward_mint.key(),
            rng_threshold: reward_params.rng_threshold,
            amount: reward_params.amount,
            cardinality: reward_params.cardinality,
        };
        quest.rewards.push(reward_params);

        Ok(())
    }

    pub fn reset_quest_rewards(
        ctx: Context<ResetQuestRewards>,
        _quest_index: u64,
        _quest_bump: u8,
    ) -> Result<(), Error> {
        let quest = &mut ctx.accounts.quest;
        quest.rewards = Vec::new();

        Ok(())
    }

    pub fn ammend_quest_with_entitlement(
        ctx: Context<AmmendQuestWithEntitlement>,
        quest_index: u64,
        _quest_bump: u8,
        entitlement: Reward,
    ) -> Result<(), Error> {
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

    pub fn start_quest<'info>(
        ctx: Context<'_, '_, '_, 'info, StartQuest<'info>>,
    ) -> Result<(), Error> {
        let now = Clock::get()?.unix_timestamp;
        let quest = &mut ctx.accounts.quest;
        let questor = &mut ctx.accounts.questor;
        let questee = &mut ctx.accounts.questee;
        let quest_account = &mut ctx.accounts.quest_acc;
        let token_program = &ctx.accounts.token_program;
        let initializer = &ctx.accounts.initializer;
        let pixelballz_token_account = &ctx.accounts.pixelballz_token_account;

        if questor.initializer != initializer.key() {
            return Err(QuestError::InvalidInitializer.into());
        }
        if questee.owner != questor.initializer {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if questor.initializer != initializer.key() {
            return Err(QuestError::SuspiciousTransaction.into());
        }
        if questee.pixelballz_mint != pixelballz_token_account.mint {
            return Err(QuestError::InvalidMint.into());
        }

        if questee.xp < quest.required_xp {
            return Err(QuestError::NotEnoughXp.into());
        }

        quest_account.start_time = now;
        quest_account.end_time = now + quest.duration;

        // TODO verify the creator address from nft metadata account
        token::transfer(
            CpiContext::new(
                token_program.to_account_info(),
                Transfer {
                    from: pixelballz_token_account.to_account_info(),
                    to: ctx.accounts.deposit_token_account.to_account_info(),
                    authority: ctx.accounts.initializer.to_account_info(),
                },
            ),
            1,
        )?;

        let mut remaining_accounts = 0;
        let ctx_remaining_accounts = ctx.remaining_accounts.to_vec();

        if quest.tender.is_some() {
            let tender = quest.tender.clone().unwrap();
            let tender_splits = quest.tender_splits.clone().unwrap();

            if tender_splits.len() == 0 {
                return Err(QuestError::SuspiciousTransaction.into());
            }
            let tender_token_account = &ctx_remaining_accounts[remaining_accounts];
            remaining_accounts += 1;

            let mut i = 0;
            while i < tender_splits.len() {
                let tender_split = tender_splits[i].clone();

                if tender_split.op_code == 0 {
                    let tendersplit_mint = &ctx_remaining_accounts[remaining_accounts];
                    remaining_accounts += 1;

                    token::burn(
                        CpiContext::new(
                            token_program.to_account_info(),
                            Burn {
                                mint: tendersplit_mint.to_account_info(),
                                from: tender_token_account.to_account_info(),
                                authority: initializer.to_account_info(),
                            },
                        ),
                        (tender.amount as f64 * ((tender_split.share as u64) as f64 / 100.0))
                            as u64,
                    )?;
                }

                if tender_split.op_code == 1 {
                    let tendersplit_token_account = &ctx_remaining_accounts[remaining_accounts];
                    remaining_accounts += 1;

                    token::transfer(
                        CpiContext::new(
                            token_program.to_account_info(),
                            Transfer {
                                from: tender_token_account.to_account_info(),
                                to: tendersplit_token_account.to_account_info(),
                                authority: initializer.to_account_info(),
                            },
                        ),
                        (tender.amount as f64 * ((tender_split.share as u64) as f64 / 100.0))
                            as u64,
                    )?;
                }

                i += 1;
            }
        }

        Ok(())
    }

    pub fn end_quest<'info>(
        ctx: Context<'_, '_, '_, 'info, EndQuest<'info>>,
        deposit_token_account_bump: u8,
        quest_bump: u8,
    ) -> Result<(), Error> {
        // TODO implement time lock enforcement
        let _now = Clock::get()?.unix_timestamp;
        let oracle = &ctx.accounts.oracle;
        let quest = &mut ctx.accounts.quest;
        let questor = &mut ctx.accounts.questor;
        let questee = &mut ctx.accounts.questee;
        let quest_questee_receipt = &mut ctx.accounts.quest_questee_receipt;
        let initializer = &ctx.accounts.initializer;
        let initializer_key = initializer.key();
        /*
        let seeds = &[
            QUEST_PDA_SEED,
            QUEST_PIXELBALLZ_SEED,
            initializer_key.as_ref(),
            quest_key.as_ref(),
        ];
        let (deposit_token_account, bump) = Pubkey::find_program_address(seeds, &crate::ID);
        if deposit_token_account != ctx.accounts.deposit_token_account.key() {
            return Err(QuestError::InvalidInitializer.into());
        }
        if deposit_token_account_bump != bump {
            return Err(QuestError::InvalidInitializer.into());
        }
        */
        let quest_key = quest.key();
        let questee_key = questee.key();
        let deposit_token_account_bump_bytes = deposit_token_account_bump.to_le_bytes();
        let seeds_with_bump = &[
            QUEST_PDA_SEED.as_ref(),
            questee_key.as_ref(),
            quest_key.as_ref(),
            deposit_token_account_bump_bytes.as_ref(),
        ];
        let deposit_token_account = &ctx.accounts.deposit_token_account;
        let deposit_token_account_authority = &[&seeds_with_bump[..]];

        let token_program = &ctx.accounts.token_program;

        // Transfer NFT quested back
        token::transfer(
            CpiContext::new_with_signer(
                token_program.to_account_info(),
                Transfer {
                    from: ctx.accounts.deposit_token_account.to_account_info(),
                    to: ctx.accounts.pixelballz_token_account.to_account_info(),
                    authority: deposit_token_account.to_account_info(),
                },
                deposit_token_account_authority,
            ),
            1,
        )?;

        // Attribute XP to Questor and Questee
        questee.xp += quest.xp;
        questor.xp += quest.xp;

        // RNG a number off recent blockhash
        let recent_slot_hash = &ctx.accounts.slot_hashes.data.borrow();
        let most_recent = &recent_slot_hash[12..25];
        // nominate for r/shittyprogramming 2022 meme of the year pls
        let percentile: u8 = (u64::from_le_bytes([
            most_recent[0],
            most_recent[1],
            most_recent[2],
            most_recent[3],
            most_recent[4],
            most_recent[5],
            most_recent[6],
            most_recent[7],
        ]) % 100) as u8;

        let (reward_index, reward) =
            convict_percentile_against_thresholds(percentile, quest.rewards.clone()).unwrap();

        let mut remaining_accounts: usize = 0;
        // select correct bump from
        let reward_mint = &ctx.remaining_accounts.to_vec()
            [remaining_accounts..remaining_accounts + quest.rewards.len()][reward_index as usize];
        remaining_accounts += quest.rewards.len();

        let reward_token_account = &ctx.remaining_accounts.to_vec()
            [remaining_accounts..remaining_accounts + quest.rewards.len()][reward_index as usize];
        // remaining_accounts += quest.rewards.len();

        let system_program = &ctx.accounts.system_program;
        if reward_token_account.data.borrow().len() == 0 {
            let _may_fail = associated_token::create(CpiContext::new(
                ctx.accounts.associated_token_program.to_account_info(),
                Create {
                    payer: initializer.to_account_info(),
                    associated_token: reward_token_account.to_account_info(),
                    authority: initializer.to_account_info(),
                    mint: reward_mint.to_account_info(),
                    system_program: system_program.to_account_info(),
                    token_program: token_program.to_account_info(),
                    rent: ctx.accounts.rent.to_account_info(),
                },
            ));
        }

        let quest_index = quest.index;
        let oracle_key = oracle.key();
        let quest_bump_bytes = quest_bump.to_le_bytes();
        let reward_authority_seeds = &[
            QUEST_ORACLE_SEED.as_ref(),
            oracle_key.as_ref(),
            &quest_index.to_le_bytes(),
            quest_bump_bytes.as_ref(),
        ];
        let reward_authority = &[&reward_authority_seeds[..]];

        // Mint reward token amount to initializer
        token::mint_to(
            CpiContext::new_with_signer(
                token_program.to_account_info(),
                MintTo {
                    mint: reward_mint.to_account_info(),
                    to: reward_token_account.to_account_info(),
                    authority: quest.to_account_info(),
                },
                reward_authority,
            ),
            reward.amount,
        )?;

        // TIL the number that is subtracting from a minuend is called a subtrahend

        quest_questee_receipt.owner = initializer_key;
        quest_questee_receipt.pixelballz_mint = deposit_token_account.mint;
        quest_questee_receipt.reward_mint = reward_mint.key();

        Ok(())
    }
}
