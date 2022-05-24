use crate::constants::*;
use crate::errors::*;
use crate::state::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use arrayref::array_ref;
use mpl_token_metadata::state::Metadata;
use mpl_token_metadata::state::{MAX_NAME_LENGTH, MAX_URI_LENGTH};
use std::cell::RefMut;
use std::result::Result;
use std::str::FromStr;

pub fn get_config_count(data: &RefMut<&mut [u8]>) -> Result<usize, Error> {
    return Ok(u32::from_le_bytes(*array_ref![data, CONFIG_ARRAY_START, 4]) as usize);
}

pub fn get_config_count_ref(data: &RefMut<&[u8]>) -> Result<usize, Error> {
    return Ok(u32::from_le_bytes(*array_ref![data, CONFIG_ARRAY_START, 4]) as usize);
}

pub fn get_space_for_batch(data: CandyMachineData) -> Result<usize, Error> {
    let num = CONFIG_ARRAY_START
        + 32
        + 32
        + 4
        + (data.items_available as usize) * CONFIG_LINE_SIZE
        + 8
        + 2 * ((data
            .items_available
            .checked_div(8)
            .ok_or(QuestError::NumericalOverflowError)?
            + 1) as usize);

    Ok(num)
}

pub fn get_good_index(
    arr: &mut RefMut<&mut [u8]>,
    items_available: usize,
    index: usize,
    pos: bool,
) -> Result<(usize, bool), Error> {
    let mut index_to_use = index;
    let mut taken = 1;
    let mut found = false;
    let bit_mask_vec_start = CONFIG_ARRAY_START
        + 4
        + (items_available) * CONFIG_LINE_SIZE
        + 4
        + items_available
            .checked_div(8)
            .ok_or(QuestError::NumericalOverflowError)?
        + 4;

    while taken > 0 && index_to_use < items_available {
        let my_position_in_vec = bit_mask_vec_start
            + index_to_use
                .checked_div(8)
                .ok_or(QuestError::NumericalOverflowError)?;
        /*msg!(
            "My position is {} and value there is {}",
            my_position_in_vec,
            arr[my_position_in_vec]
        );*/
        if arr[my_position_in_vec] == 255 {
            //msg!("We are screwed here, move on");
            let eight_remainder = 8 - index_to_use
                .checked_rem(8)
                .ok_or(QuestError::NumericalOverflowError)?;
            let reversed = 8 - eight_remainder + 1;
            if (eight_remainder != 0 && pos) || (reversed != 0 && !pos) {
                //msg!("Moving by {}", eight_remainder);
                if pos {
                    index_to_use += eight_remainder;
                } else {
                    if index_to_use < 8 {
                        break;
                    }
                    index_to_use -= reversed;
                }
            } else {
                //msg!("Moving by 8");
                if pos {
                    index_to_use += 8;
                } else {
                    index_to_use -= 8;
                }
            }
        } else {
            let position_from_right = 7 - index_to_use
                .checked_rem(8)
                .ok_or(QuestError::NumericalOverflowError)?;
            let mask = u8::pow(2, position_from_right as u32);

            taken = mask & arr[my_position_in_vec];
            if taken > 0 {
                //msg!("Index to use {} is taken", index_to_use);
                if pos {
                    index_to_use += 1;
                } else {
                    if index_to_use == 0 {
                        break;
                    }
                    index_to_use -= 1;
                }
            } else if taken == 0 {
                //msg!("Index to use {} is not taken, exiting", index_to_use);
                found = true;
                arr[my_position_in_vec] = arr[my_position_in_vec] | mask;
            }
        }
    }

    Ok((index_to_use, found))
}

pub fn get_config_line<'info>(
    a: &Account<'info, Batch>,
    index_to_use: usize,
) -> Result<ConfigLine, Error> {
    msg!("Index is set to {:?}", index_to_use);
    let a_info = a.to_account_info();
    let mut arr = a_info.data.borrow_mut();

    msg!(
        "Index actually ends up due to used bools {:?}",
        index_to_use
    );
    if arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)] == 1 {
        return Err(QuestError::CannotFindUsableConfigLine.into());
    }

    let data_array = &mut arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)
        ..CONFIG_ARRAY_START + 4 + (index_to_use + 1) * (CONFIG_LINE_SIZE)];

    let mut name_vec = vec![];
    let mut cardinality_vec = vec![];
    let mut uri_vec = vec![];
    for i in 4..4 + MAX_NAME_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        name_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH..8 + MAX_NAME_LENGTH + MAX_CARDINALITY_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        cardinality_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH + 4 + MAX_CARDINALITY_LENGTH
        ..8 + MAX_CARDINALITY_LENGTH + MAX_NAME_LENGTH + MAX_URI_LENGTH
    {
        if data_array[i] == 0 {
            break;
        }
        uri_vec.push(data_array[i])
    }
    let config_line: ConfigLine = ConfigLine {
        name: match String::from_utf8(name_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        cardinality: match String::from_utf8(cardinality_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        uri: match String::from_utf8(uri_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
    };

    Ok(config_line)
}
pub fn get_config_line_from_account_info<'info>(
    a_info: &AccountInfo,
    index_to_use: usize,
) -> Result<ConfigLine, Error> {
    let mut arr = a_info.data.borrow_mut();

    if arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)] == 1 {
        return Err(QuestError::CannotFindUsableConfigLine.into());
    }

    let data_array = &mut arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)
        ..CONFIG_ARRAY_START + 4 + (index_to_use + 1) * (CONFIG_LINE_SIZE)];

    let mut name_vec = vec![];
    let mut cardinality_vec = vec![];
    let mut uri_vec = vec![];
    for i in 4..4 + MAX_NAME_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        name_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH..8 + MAX_NAME_LENGTH + MAX_CARDINALITY_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        cardinality_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH + 4 + MAX_CARDINALITY_LENGTH
        ..8 + MAX_CARDINALITY_LENGTH + MAX_NAME_LENGTH + MAX_URI_LENGTH
    {
        if data_array[i] == 0 {
            break;
        }
        uri_vec.push(data_array[i])
    }
    let config_line: ConfigLine = ConfigLine {
        name: match String::from_utf8(name_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        cardinality: match String::from_utf8(cardinality_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        uri: match String::from_utf8(uri_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
    };

    Ok(config_line)
}

pub fn get_config_lines<'info>(
    arr: RefMut<&[u8]>,
    index_to_use: usize,
) -> Result<ConfigLine, Error> {
    if arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)] == 1 {
        return Err(QuestError::CannotFindUsableConfigLine.into());
    }

    let data_array = &arr[CONFIG_ARRAY_START + 4 + index_to_use * (CONFIG_LINE_SIZE)
        ..CONFIG_ARRAY_START + 4 + (index_to_use + 1) * (CONFIG_LINE_SIZE)];

    let mut name_vec = vec![];
    let mut cardinality_vec = vec![];
    let mut uri_vec = vec![];
    for i in 4..4 + MAX_NAME_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        name_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH..8 + MAX_NAME_LENGTH + MAX_CARDINALITY_LENGTH {
        if data_array[i] == 0 {
            break;
        }
        cardinality_vec.push(data_array[i])
    }
    for i in 8 + MAX_NAME_LENGTH + 4 + MAX_CARDINALITY_LENGTH
        ..8 + MAX_CARDINALITY_LENGTH + MAX_NAME_LENGTH + MAX_URI_LENGTH
    {
        if data_array[i] == 0 {
            break;
        }
        uri_vec.push(data_array[i])
    }
    let config_line: ConfigLine = ConfigLine {
        name: match String::from_utf8(name_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        cardinality: match String::from_utf8(cardinality_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
        uri: match String::from_utf8(uri_vec) {
            Ok(val) => val,
            Err(_) => return Err(QuestError::InvalidString.into()),
        },
    };

    Ok(config_line)
}

pub fn assert_valid_metadata(
    depo_metadata: &AccountInfo,
    depo_mint: &Pubkey,
) -> Result<Metadata, ProgramError> {
    let metadata_program = Pubkey::from_str("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s").unwrap();
    assert_eq!(depo_metadata.owner, &metadata_program);
    let seed = &[
        b"metadata".as_ref(),
        metadata_program.as_ref(),
        depo_mint.as_ref(),
    ];
    let (metadata_addr, _bump) = Pubkey::find_program_address(seed, &metadata_program);
    assert_eq!(metadata_addr, depo_metadata.key());

    Metadata::from_account_info(depo_metadata)
}

pub fn assert_all_unique_account_infos(account_infos: &Vec<AccountInfo>) -> Result<(), Error> {
    let account_infos_cloned = &account_infos.clone();

    let mut i: usize = 0;
    while i < account_infos_cloned.len() {
        let account_info = account_infos_cloned[i].clone();
        let account_info_key = account_info.key();
        let dupe_account_infos = account_infos_cloned
            .into_iter()
            .filter(|n| n.key() == account_info_key)
            .collect::<Vec<_>>();

        if dupe_account_infos.len() != 1 {
            return Err(QuestError::SuspiciousAccounts.into());
        }
        i += 1;
    }

    Ok(())
}

pub fn get_valid_batch_for_cardinality(
    batch_account_infos: &Vec<AccountInfo>,
    rng: u64,
    cardinality: String,
) -> Result<(Pubkey, u64), Error> {
    let batch_account_infos_cloned = batch_account_infos.clone();
    let mut result: Option<(Pubkey, u64)> = None;
    let max_depth: u8 = batch_account_infos.len() as u8;
    let mut depth: usize = 0;

    // if the retry depth exceeds number of batch accounts,
    // then panic after searching as we have exhausted
    // every batch account for a cardinality
    // and it yielded nothing.
    while depth < max_depth as usize {
        // rng_batch_indice is a u8 since it represents an batch account index
        // from a batch accounts array proposed in an instruction. so it is safe
        // to perform math as u8 types.
        let supposed_rng_batch_indice = ((rng + depth as u64) % max_depth as u64) as usize;
        let batch_cardinalities_report = BatchCardinalitiesReport::from_account_info(
            &batch_account_infos_cloned[supposed_rng_batch_indice],
        )?;

        let cardinality_indice = batch_cardinalities_report
            .cardinalities_keys
            .into_iter()
            .position(|n| n == cardinality);

        if cardinality_indice.is_some() {
            if batch_cardinalities_report.cardinalities_indices[cardinality_indice.unwrap()].len()
                > 0
            {
                result = Some((
                    batch_cardinalities_report.batch_account,
                    cardinality_indice.unwrap() as u64,
                ));
                break;
            }
        }

        depth += 1;
    }

    // `.unwrap()` is to be panicked if `None` should
    // the above search fail to set an actual tuple.
    Ok(result.unwrap())
}
