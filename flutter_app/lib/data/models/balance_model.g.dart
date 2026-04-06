// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'balance_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

TotalBalanceModel _$TotalBalanceModelFromJson(Map<String, dynamic> json) =>
    TotalBalanceModel(
      details: json['details'] == null
          ? null
          : BalanceDetails.fromJson(json['details'] as Map<String, dynamic>),
      total: json['total'] == null
          ? null
          : BalanceTotal.fromJson(json['total'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$TotalBalanceModelToJson(TotalBalanceModel instance) =>
    <String, dynamic>{'details': instance.details, 'total': instance.total};

BalanceDetails _$BalanceDetailsFromJson(Map<String, dynamic> json) =>
    BalanceDetails(
      crossMargin: json['cross_margin'] == null
          ? null
          : BalanceItem.fromJson(json['cross_margin'] as Map<String, dynamic>),
      spot: json['spot'] == null
          ? null
          : BalanceItem.fromJson(json['spot'] as Map<String, dynamic>),
      finance: json['finance'] == null
          ? null
          : BalanceItem.fromJson(json['finance'] as Map<String, dynamic>),
      margin: json['margin'] == null
          ? null
          : BalanceItem.fromJson(json['margin'] as Map<String, dynamic>),
      quant: json['quant'] == null
          ? null
          : BalanceItem.fromJson(json['quant'] as Map<String, dynamic>),
      futures: json['futures'] == null
          ? null
          : BalanceItem.fromJson(json['futures'] as Map<String, dynamic>),
      delivery: json['delivery'] == null
          ? null
          : BalanceItem.fromJson(json['delivery'] as Map<String, dynamic>),
      warrant: json['warrant'] == null
          ? null
          : BalanceItem.fromJson(json['warrant'] as Map<String, dynamic>),
      cbbc: json['cbbc'] == null
          ? null
          : BalanceItem.fromJson(json['cbbc'] as Map<String, dynamic>),
      memeBox: json['memeBox'] == null
          ? null
          : BalanceItem.fromJson(json['memeBox'] as Map<String, dynamic>),
      options: json['options'] == null
          ? null
          : BalanceItem.fromJson(json['options'] as Map<String, dynamic>),
      payment: json['payment'] == null
          ? null
          : BalanceItem.fromJson(json['payment'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$BalanceDetailsToJson(BalanceDetails instance) =>
    <String, dynamic>{
      'cross_margin': instance.crossMargin,
      'spot': instance.spot,
      'finance': instance.finance,
      'margin': instance.margin,
      'quant': instance.quant,
      'futures': instance.futures,
      'delivery': instance.delivery,
      'warrant': instance.warrant,
      'cbbc': instance.cbbc,
      'memeBox': instance.memeBox,
      'options': instance.options,
      'payment': instance.payment,
    };

BalanceItem _$BalanceItemFromJson(Map<String, dynamic> json) => BalanceItem(
  amount: json['amount'] as String,
  currency: json['currency'] as String,
  borrowed: json['borrowed'] as String?,
  unrealisedPnl: json['unrealised_pnl'] as String?,
);

Map<String, dynamic> _$BalanceItemToJson(BalanceItem instance) =>
    <String, dynamic>{
      'amount': instance.amount,
      'currency': instance.currency,
      'borrowed': instance.borrowed,
      'unrealised_pnl': instance.unrealisedPnl,
    };

BalanceTotal _$BalanceTotalFromJson(Map<String, dynamic> json) => BalanceTotal(
  amount: json['amount'] as String,
  currency: json['currency'] as String,
  unrealisedPnl: json['unrealised_pnl'] as String?,
  borrowed: json['borrowed'] as String?,
);

Map<String, dynamic> _$BalanceTotalToJson(BalanceTotal instance) =>
    <String, dynamic>{
      'amount': instance.amount,
      'currency': instance.currency,
      'unrealised_pnl': instance.unrealisedPnl,
      'borrowed': instance.borrowed,
    };
