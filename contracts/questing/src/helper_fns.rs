use crate::errors::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use std::result::Result;

pub fn assert_valid_sum_rng_thresholds(
    rewards: Vec<Reward>,
    additional_rng_threshold: u8,
) -> Result<(), Error> {
    let mut sum: u8 = additional_rng_threshold;
    let mut i = 0;
    while i < rewards.len() {
        if sum > 100 {
            msg!("sum rng theshold must not exceed 100!");
            return Err(QuestError::IndexGreaterThanLength.into());
        }
        sum += rewards[i].rng_threshold;
        i += 1;
    }

    Ok(())
}

pub fn get_sum_xp_from_level(level: u64) -> Result<u64, Error> {
    if level == 0 {
        return Ok(level);
    }
    let base_xp: u64 = 50;
    return Ok((u64::pow(2, level as u32)) * base_xp);
}

pub fn convict_percentile_against_thresholds(
    percentile: u8,
    thresholds: Vec<Reward>,
) -> Result<(u64, Reward), Error> {
    let mut accumulator: u8 = 0;
    let mut i: u64 = 0;
    let thresholds_length = thresholds.len() as u64;
    while i < thresholds_length {
        let pct = thresholds[i as usize].clone();
        let previous_accumulator = accumulator;
        accumulator += pct.rng_threshold;
        if previous_accumulator <= percentile && percentile <= accumulator {
            return Ok((i, pct));
        }
        i += 1;
    }

    Err(QuestError::InvalidConviction.into())
}
