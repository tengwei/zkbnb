/*
 * Copyright © 2021 Zecrey Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package proverUtil

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	curve "github.com/zecrey-labs/zecrey-crypto/ecc/ztwistededwards/tebn254"
	"github.com/zecrey-labs/zecrey-crypto/zecrey-legend/circuit/bn254/std"
	"github.com/zecrey-labs/zecrey-legend/common/commonTx"
	"github.com/zeromicro/go-zero/core/logx"
)

func ConstructRegisterZnsCryptoTx(
	oTx *Tx,
	accountTree *Tree,
	accountAssetsTree *[]*Tree,
	liquidityTree *Tree,
	nftTree *Tree,
	accountModel AccountModel,
) (cryptoTx *CryptoTx, err error) {
	if oTx == nil || accountTree == nil || accountAssetsTree == nil || liquidityTree == nil || nftTree == nil {
		logx.Errorf("[ConstructRegisterZnsCryptoTx] invalid params")
		return nil, errors.New("[ConstructRegisterZnsCryptoTx] invalid params")
	}
	txInfo, err := commonTx.ParseRegisterZnsTxInfo(oTx.TxInfo)
	if err != nil {
		logx.Errorf("[ConstructRegisterZnsCryptoTx] unable to parse register zns tx info:%s", err.Error())
		return nil, err
	}
	cryptoTxInfo, err := ToCryptoRegisterZnsTx(txInfo)
	if err != nil {
		logx.Errorf("[ConstructRegisterZnsCryptoTx] unable to convert to crypto register zns tx: %s", err.Error())
		return nil, err
	}
	accountKeys, proverAccountMap, proverLiquidityInfo, proverNftInfo, err := ConstructProverInfo(oTx, accountModel)
	if err != nil {
		logx.Errorf("[ConstructRegisterZnsCryptoTx] unable to construct prover info: %s", err.Error())
		return nil, err
	}
	cryptoTx, err = ConstructWitnessInfo(
		oTx,
		accountTree,
		accountAssetsTree,
		liquidityTree,
		nftTree,
		accountKeys,
		proverAccountMap,
		proverLiquidityInfo,
		proverNftInfo,
	)
	if err != nil {
		logx.Errorf("[ConstructRegisterZnsCryptoTx] unable to construct witness info: %s", err.Error())
		return nil, err
	}
	cryptoTx.RegisterZnsTxInfo = cryptoTxInfo
	cryptoTx.Nonce = oTx.Nonce
	cryptoTx.Signature = std.EmptySignature()
	return cryptoTx, nil
}

func ToCryptoRegisterZnsTx(txInfo *commonTx.RegisterZnsTxInfo) (info *CryptoRegisterZnsTx, err error) {
	accountName := make([]byte, 32)
	copy(accountName[:], txInfo.AccountName)
	pubKeyBytes, err := hex.DecodeString(txInfo.PubKey)
	if err != nil {
		logx.Errorf("[ToCryptoRegisterZnsTx] unable to decode string:%s", err.Error())
		return nil, err
	}
	pk, err := curve.FromBytes(pubKeyBytes)
	if err != nil {
		logx.Errorf("[ToCryptoRegisterZnsTx] unable to parse pub key:%s", err.Error())
		return nil, err
	}
	info = &CryptoRegisterZnsTx{
		AccountIndex:    txInfo.AccountIndex,
		AccountName:     accountName,
		AccountNameHash: common.FromHex(txInfo.AccountNameHash),
	}
	info.PubKey.A.Set(pk)
	return info, nil
}
