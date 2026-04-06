// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'futures_account_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FuturesAccountModel _$FuturesAccountModelFromJson(Map<String, dynamic> json) =>
    FuturesAccountModel(
      total: json['total'] as String,
      unrealisedPnl: json['unrealisedPnl'] as String,
      available: json['available'] as String,
      orderMargin: json['orderMargin'] as String,
      positionMargin: json['positionMargin'] as String,
      maintenanceMargin: json['maintenanceMargin'] as String,
      currency: json['currency'] as String,
      inDualMode: json['inDualMode'] as bool,
    );

Map<String, dynamic> _$FuturesAccountModelToJson(
  FuturesAccountModel instance,
) => <String, dynamic>{
  'total': instance.total,
  'unrealisedPnl': instance.unrealisedPnl,
  'available': instance.available,
  'orderMargin': instance.orderMargin,
  'positionMargin': instance.positionMargin,
  'maintenanceMargin': instance.maintenanceMargin,
  'currency': instance.currency,
  'inDualMode': instance.inDualMode,
};
