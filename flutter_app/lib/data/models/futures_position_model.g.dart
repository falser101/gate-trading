// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'futures_position_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FuturesPositionModel _$FuturesPositionModelFromJson(
  Map<String, dynamic> json,
) => FuturesPositionModel(
  contract: json['contract'] as String?,
  size: json['size'] as String?,
  leverage: json['leverage'] as String?,
  entryPrice: json['entryPrice'] as String?,
  markPrice: json['markPrice'] as String?,
  liqPrice: json['liqPrice'] as String?,
  margin: json['margin'] as String?,
  unrealisedPnl: json['unrealisedPnl'] as String?,
  realisedPnl: json['realisedPnl'] as String?,
  initialMargin: json['initialMargin'] as String?,
  maintenanceMargin: json['maintenanceMargin'] as String?,
  adlRanking: (json['adlRanking'] as num?)?.toInt(),
  mode: json['mode'] as String?,
);

Map<String, dynamic> _$FuturesPositionModelToJson(
  FuturesPositionModel instance,
) => <String, dynamic>{
  'contract': instance.contract,
  'size': instance.size,
  'leverage': instance.leverage,
  'entryPrice': instance.entryPrice,
  'markPrice': instance.markPrice,
  'liqPrice': instance.liqPrice,
  'margin': instance.margin,
  'unrealisedPnl': instance.unrealisedPnl,
  'realisedPnl': instance.realisedPnl,
  'initialMargin': instance.initialMargin,
  'maintenanceMargin': instance.maintenanceMargin,
  'adlRanking': instance.adlRanking,
  'mode': instance.mode,
};
