// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'account_detail_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AccountDetailModel _$AccountDetailModelFromJson(Map<String, dynamic> json) =>
    AccountDetailModel(
      userId: (json['userId'] as num?)?.toInt(),
      ipWhitelist: (json['ip_whitelist'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      currencyPairs: (json['currency_pairs'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      key: json['key'] == null
          ? null
          : AccountKeyInfo.fromJson(json['key'] as Map<String, dynamic>),
      tier: (json['tier'] as num?)?.toInt(),
      copyTradingRole: (json['copy_trading_role'] as num?)?.toInt(),
    );

Map<String, dynamic> _$AccountDetailModelToJson(AccountDetailModel instance) =>
    <String, dynamic>{
      'userId': instance.userId,
      'ip_whitelist': instance.ipWhitelist,
      'currency_pairs': instance.currencyPairs,
      'key': instance.key,
      'tier': instance.tier,
      'copy_trading_role': instance.copyTradingRole,
    };

AccountKeyInfo _$AccountKeyInfoFromJson(Map<String, dynamic> json) =>
    AccountKeyInfo(mode: (json['mode'] as num?)?.toInt());

Map<String, dynamic> _$AccountKeyInfoToJson(AccountKeyInfo instance) =>
    <String, dynamic>{'mode': instance.mode};
