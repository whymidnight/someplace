use crate::constants::*;
use crate::errors::*;
use crate::state::*;
use crate::structs::*;
use anchor_lang::prelude::*;
use arrayref::array_ref;
use mpl_token_metadata::state::Metadata;
use mpl_token_metadata::state::{MAX_NAME_LENGTH, MAX_URI_LENGTH};
use std::cell::RefMut;
use std::str::FromStr;

pub fn get_config_count(data: &RefMut<&mut [u8]>) -> core::result::Result<usize, ProgramError> {
    return Ok(u32::from_le_bytes(*array_ref![data, CONFIG_ARRAY_START, 4]) as usize);
}

pub fn get_config_count_ref(data: &RefMut<&[u8]>) -> core::result::Result<usize, ProgramError> {
    return Ok(u32::from_le_bytes(*array_ref![data, CONFIG_ARRAY_START, 4]) as usize);
}

pub fn get_space_for_batch(data: CandyMachineData) -> core::result::Result<usize, ProgramError> {
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
) -> core::result::Result<(usize, bool), ProgramError> {
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
) -> core::result::Result<ConfigLine, ProgramError> {
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

pub fn get_config_lines<'info>(
    arr: RefMut<&[u8]>,
    index_to_use: usize,
) -> core::result::Result<ConfigLine, ProgramError> {
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
) -> core::result::Result<Metadata, ProgramError> {
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

