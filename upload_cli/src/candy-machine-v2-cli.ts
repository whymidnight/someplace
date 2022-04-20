//@ts-nocheck
import * as fs from 'fs';
import * as path from 'path';
import { InvalidArgumentError, program } from 'commander';
import * as anchor from '@project-serum/anchor';

import {
  chunks,
  fromUTF8Array,
  getCandyMachineV2Config,
  parsePrice,
} from './helpers/various';
import { PublicKey, LAMPORTS_PER_SOL } from '@solana/web3.js';
import {
  CACHE_PATH,
  CONFIG_LINE_SIZE_V2,
  EXTENSION_JSON,
  CANDY_MACHINE_PROGRAM_V2_ID,
  CONFIG_ARRAY_START_V2,
} from './helpers/constants';
import {
  getProgramAccounts,
  loadCandyProgramV2,
  loadWalletKey,
  AccountAndPubkey,
  deriveCandyMachineV2ProgramAddress,
} from './helpers/accounts';

import { uploadV2 } from './commands/upload';
import { verifyTokenMetadata } from './commands/verifyTokenMetadata';
import { loadCache, saveCache } from './helpers/cache';
import { mintV2 } from './commands/mint';
import { signMetadata } from './commands/sign';
import {
  getAccountsByCreatorAddress,
  signAllMetadataFromCandyMachine,
} from './commands/signAll';
import log from 'loglevel';
import { withdrawV2 } from './commands/withdraw';
import { updateFromCache } from './commands/updateFromCache';
import { StorageType } from './helpers/storage-type';
import { getType } from 'mime';
program.version('0.0.2');
const supportedImageTypes = {
  'image/png': 1,
  'image/gif': 1,
  'image/jpeg': 1,
};
const supportedAnimationTypes = {
  'video/mp4': 1,
  'video/quicktime': 1,
  'audio/mpeg': 1,
  'audio/x-flac': 1,
  'audio/wav': 1,
  'model/gltf-binary': 1,
  'text/html': 1,
};

if (!fs.existsSync(CACHE_PATH)) {
  fs.mkdirSync(CACHE_PATH);
}
log.setLevel(log.levels.INFO);

// From commander examples
function myParseInt(value) {
  // parseInt takes a string and a radix
  const parsedValue = parseInt(value, 10);
  if (isNaN(parsedValue)) {
    throw new InvalidArgumentError('Not a number.');
  }
  return parsedValue;
}

programCommand('upload')
  .argument(
    '<directory>',
    'Directory containing images named from 0-n',
    val => {
      return fs.readdirSync(`${val}`).map(file => path.join(val, file));
    },
  )
  .requiredOption(
    '-cp, --config-path <string>',
    'JSON file with candy machine settings',
  )
  .option(
    '-r, --rpc-url <string>',
    'custom rpc url since this is a heavy command',
  )
  .option(
    '-rl, --rate-limit <number>',
    'max number of requests per second',
    myParseInt,
    5,
  )
  .action(async (files: string[], options, cmd) => {
    const { keypair, env, cacheName, configPath, rpcUrl, rateLimit } =
      cmd.opts();

    const walletKeyPair = loadWalletKey(keypair);
    const anchorProgram = await loadCandyProgramV2(walletKeyPair, env, rpcUrl);

    const {
      storage,
      nftStorageKey,
      ipfsInfuraProjectId,
      number,
      ipfsInfuraSecret,
      arweaveJwk,
      awsS3Bucket,
      retainAuthority,
      mutable,
      batchSize,
      price,
      splToken,
      treasuryWallet,
      gatekeeper,
      endSettings,
      hiddenSettings,
      whitelistMintSettings,
      goLiveDate,
      uuid,
    } = await getCandyMachineV2Config(walletKeyPair, anchorProgram, configPath);

    console.log(".,......,,");
    if (storage === StorageType.ArweaveSol && env !== 'mainnet-beta') {
      log.info(
        'The arweave-sol storage option only works on mainnet. For devnet, please use either arweave, aws or ipfs\n',
      );
    }

    if (storage === StorageType.ArweaveBundle && env !== 'mainnet-beta') {
      throw new Error(
        'The arweave-bundle storage option only works on mainnet because it requires spending real AR tokens. For devnet, please set the --storage option to "aws" or "ipfs"\n',
      );
    }

    if (storage === StorageType.Arweave) {
      log.warn(
        'WARNING: The "arweave" storage option will be going away soon. Please migrate to arweave-bundle or arweave-sol for mainnet.\n',
      );
    }

    if (storage === StorageType.ArweaveBundle && !arweaveJwk) {
      throw new Error(
        'Path to Arweave JWK wallet file (--arweave-jwk) must be provided when using arweave-bundle',
      );
    }
    if (
      storage === StorageType.Ipfs &&
      (!ipfsInfuraProjectId || !ipfsInfuraSecret)
    ) {
      throw new Error(
        'IPFS selected as storage option but Infura project id or secret key were not provided.',
      );
    }
    if (storage === StorageType.Aws && !awsS3Bucket) {
      throw new Error(
        'aws selected as storage option but existing bucket name (--aws-s3-bucket) not provided.',
      );
    }
    if (!Object.values(StorageType).includes(storage)) {
      throw new Error(
        `Storage option must either be ${Object.values(StorageType).join(
          ', ',
        )}. Got: ${storage}`,
      );
    }
    const ipfsCredentials = {
      projectId: ipfsInfuraProjectId,
      secretKey: ipfsInfuraSecret,
    };

    let imageFileCount = 0;
    let animationFileCount = 0;
    let jsonFileCount = 0;

    // Filter out any non-supported file types and find the JSON vs Image file count
    const supportedFiles = files.filter(it => {
      if (supportedImageTypes[getType(it)]) {
        imageFileCount++;
      } else if (supportedAnimationTypes[getType(it)]) {
        animationFileCount++;
      } else if (it.endsWith(EXTENSION_JSON)) {
        jsonFileCount++;
      } else {
        log.warn(`WARNING: Skipping unsupported file type ${it}`);
        return false;
      }

      return true;
    });

    if (animationFileCount !== 0 && storage === StorageType.Arweave) {
      throw new Error(
        'The "arweave" storage option is incompatible with animation files. Please try again with another storage option using `--storage <option>`.',
      );
    }

    if (animationFileCount !== 0 && animationFileCount !== jsonFileCount) {
      throw new Error(
        `number of animation files (${animationFileCount}) is different than the number of json files (${jsonFileCount})`,
      );
    } else if (imageFileCount !== jsonFileCount) {
      throw new Error(
        `number of img files (${imageFileCount}) is different than the number of json files (${jsonFileCount})`,
      );
    }

    const elemCount = number ? number : imageFileCount;
    if (elemCount < imageFileCount) {
      throw new Error(
        `max number (${elemCount}) cannot be smaller than the number of images in the source folder (${imageFileCount})`,
      );
    }

    if (animationFileCount === 0) {
      log.info(`Beginning the upload for ${elemCount} (img+json) pairs`);
    } else {
      log.info(
        `Beginning the upload for ${elemCount} (img+animation+json) sets`,
      );
    }

    const startMs = Date.now();
    log.info('started at: ' + startMs.toString());
    console.log(".,......,,");
    try {
      console.log(".,......,,");
      await uploadV2({
        files: supportedFiles,
        cacheName,
        env,
        totalNFTs: elemCount,
        gatekeeper,
        storage,
        retainAuthority,
        mutable,
        nftStorageKey,
        ipfsCredentials,
        awsS3Bucket,
        batchSize,
        price,
        treasuryWallet,
        anchorProgram,
        walletKeyPair,
        splToken,
        endSettings,
        hiddenSettings,
        whitelistMintSettings,
        goLiveDate,
        uuid,
        arweaveJwk,
        rateLimit,
      });
    } catch (err) {
      log.warn('upload was not successful, please re-run.', err);
      process.exit(1);
    }
    const endMs = Date.now();
    const timeTaken = new Date(endMs - startMs).toISOString().substr(11, 8);
    log.info(
      `ended at: ${new Date(endMs).toISOString()}. time taken: ${timeTaken}`,
    );
    process.exit(0);
  });


program
  .command('verify_assets')
  .argument(
    '<directory>',
    'Directory containing images and metadata files named from 0-n',
    val => {
      return fs
        .readdirSync(`${val}`)
        .map(file => path.join(process.cwd(), val, file));
    },
  )
  .option('-n, --number <number>', 'Number of images to upload')
  .action((files: string[], options, cmd) => {
    const { number } = cmd.opts();

    const startMs = Date.now();
    log.info('started at: ' + startMs.toString());
    verifyTokenMetadata({ files, uploadElementsCount: number });

    const endMs = Date.now();
    const timeTaken = new Date(endMs - startMs).toISOString().substr(11, 8);
    log.info(
      `ended at: ${new Date(endMs).toString()}. time taken: ${timeTaken}`,
    );
  });

programCommand('verify_upload')
  .option(
    '-r, --rpc-url <string>',
    'custom rpc url since this is a heavy command',
  )
  .action(async (directory, cmd) => {
    const { env, keypair, rpcUrl, cacheName } = cmd.opts();

    const cacheContent = loadCache(cacheName, env);
    const walletKeyPair = loadWalletKey(keypair);
    const anchorProgram = await loadCandyProgramV2(walletKeyPair, env, rpcUrl);

    const candyMachine = await anchorProgram.provider.connection.getAccountInfo(
      new PublicKey(cacheContent.program.candyMachine),
    );

    const candyMachineObj = await anchorProgram.account.candyMachine.fetch(
      new PublicKey(cacheContent.program.candyMachine),
    );
    let allGood = true;

    const keys = Object.keys(cacheContent.items)
      .filter(k => !cacheContent.items[k].verifyRun)
      .sort((a, b) => Number(a) - Number(b));

    console.log('Key size', keys.length);
    await Promise.all(
      chunks(keys, 500).map(async allIndexesInSlice => {
        for (let i = 0; i < allIndexesInSlice.length; i++) {
          // Save frequently.
          if (i % 100 == 0) saveCache(cacheName, env, cacheContent);

          const key = allIndexesInSlice[i];
          log.info('Looking at key ', key);

          const thisSlice = candyMachine.data.slice(
            CONFIG_ARRAY_START_V2 + 4 + CONFIG_LINE_SIZE_V2 * key,
            CONFIG_ARRAY_START_V2 + 4 + CONFIG_LINE_SIZE_V2 * (key + 1),
          );

          const name = fromUTF8Array([
            ...thisSlice.slice(4, 36).filter(n => n !== 0),
          ]);
          const uri = fromUTF8Array([
            ...thisSlice.slice(40, 240).filter(n => n !== 0),
          ]);
          const cacheItem = cacheContent.items[key];

          if (name != cacheItem.name || uri != cacheItem.link) {
            //leaving here for debugging reasons, but it's pretty useless. if the first upload fails - all others are wrong
            /*log.info(
                `Name (${name}) or uri (${uri}) didnt match cache values of (${cacheItem.name})` +
                  `and (${cacheItem.link}). marking to rerun for image`,
                key,
              );*/
            cacheItem.onChain = false;
            allGood = false;
          } else {
            cacheItem.verifyRun = true;
          }
        }
      }),
    );

    if (!allGood) {
      saveCache(cacheName, env, cacheContent);

      throw new Error(
        `not all NFTs checked out. check out logs above for details`,
      );
    }

    const lineCount = new anchor.BN(
      candyMachine.data.slice(CONFIG_ARRAY_START_V2, CONFIG_ARRAY_START_V2 + 4),
      undefined,
      'le',
    );

    log.info(
      `uploaded (${lineCount.toNumber()}) out of (${
        candyMachineObj.data.itemsAvailable
      })`,
    );
    if (candyMachineObj.data.itemsAvailable > lineCount.toNumber()) {
      throw new Error(
        `predefined number of NFTs (${
          candyMachineObj.data.itemsAvailable
        }) is smaller than the uploaded one (${lineCount.toNumber()})`,
      );
    } else {
      log.info('ready to deploy!');
    }

    saveCache(cacheName, env, cacheContent);
  });

programCommand('sign')
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  .requiredOption('-m, --metadata <string>', 'base58 metadata account id')
  .option(
    '-r, --rpc-url <string>',
    'custom rpc url since this is a heavy command',
  )
  .action(async (directory, cmd) => {
    const { keypair, env, rpcUrl, metadata } = cmd.opts();

    await signMetadata(metadata, keypair, env, rpcUrl);
  });

programCommand('sign_all')
  .option('-b, --batch-size <string>', 'Batch size', '10')
  .option('-d, --daemon', 'Run signing continuously', false)
  .option(
    '-r, --rpc-url <string>',
    'custom rpc url since this is a heavy command',
  )
  .action(async (directory, cmd) => {
    const { keypair, env, cacheName, rpcUrl, batchSize, daemon } = cmd.opts();
    const cacheContent = loadCache(cacheName, env);
    const walletKeyPair = loadWalletKey(keypair);
    const anchorProgram = await loadCandyProgramV2(walletKeyPair, env, rpcUrl);

    const batchSizeParsed = parseInt(batchSize);
    if (!parseInt(batchSize)) {
      throw new Error('Batch size needs to be an integer!');
    }

    const candyMachineId = new PublicKey(cacheContent.program.candyMachine);
    const [candyMachineAddr] = await deriveCandyMachineV2ProgramAddress(
      candyMachineId,
    );

    log.debug('Creator pubkey: ', walletKeyPair.publicKey.toBase58());
    log.debug('Environment: ', env);
    log.debug('Candy machine address: ', cacheContent.program.candyMachine);
    log.debug('Batch Size: ', batchSizeParsed);
    await signAllMetadataFromCandyMachine(
      anchorProgram.provider.connection,
      walletKeyPair,
      candyMachineAddr.toBase58(),
      batchSizeParsed,
      daemon,
    );
  });

programCommand('update_existing_nfts_from_latest_cache_file')
  .option('-b, --batch-size <string>', 'Batch size', '2')
  .option('-nc, --new-cache <string>', 'Path to new updated cache file')
  .option('-d, --daemon', 'Run updating continuously', false)
  .option(
    '-r, --rpc-url <string>',
    'custom rpc url since this is a heavy command',
  )
  .action(async (directory, cmd) => {
    const { keypair, env, cacheName, rpcUrl, batchSize, daemon, newCache } =
      cmd.opts();
    const cacheContent = loadCache(cacheName, env);
    const newCacheContent = loadCache(newCache, env);
    const walletKeyPair = loadWalletKey(keypair);
    const anchorProgram = await loadCandyProgramV2(walletKeyPair, env, rpcUrl);
    const candyAddress = cacheContent.program.candyMachine;

    const batchSizeParsed = parseInt(batchSize);
    if (!parseInt(batchSize)) {
      throw new Error('Batch size needs to be an integer!');
    }

    log.debug('Creator pubkey: ', walletKeyPair.publicKey.toBase58());
    log.debug('Environment: ', env);
    log.debug('Candy machine address: ', candyAddress);
    log.debug('Batch Size: ', batchSizeParsed);
    await updateFromCache(
      anchorProgram.provider.connection,
      walletKeyPair,
      candyAddress,
      batchSizeParsed,
      daemon,
      cacheContent,
      newCacheContent,
    );
  });


function programCommand(name: string) {
  return program
    .command(name)
    .option(
      '-e, --env <string>',
      'Solana cluster env name',
      'devnet', //mainnet-beta, testnet, devnet
    )
    .requiredOption('-k, --keypair <path>', `Solana wallet location`)
    .option('-l, --log-level <string>', 'log level', setLogLevel)
    .option('-c, --cache-name <string>', 'Cache file name', 'temp');
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function setLogLevel(value, prev) {
  if (value === undefined || value === null) {
    return;
  }
  log.info('setting the log value to: ' + value);
  log.setLevel(value);
}

program.parse(process.argv);
