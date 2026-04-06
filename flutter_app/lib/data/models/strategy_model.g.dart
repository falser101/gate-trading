// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'strategy_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StrategyModel _$StrategyModelFromJson(Map<String, dynamic> json) =>
    StrategyModel(
      id: (json['id'] as num).toInt(),
      userId: (json['user_id'] as num).toInt(),
      name: json['name'] as String?,
      type: json['type'] as String,
      symbol: json['symbol'] as String,
      config: json['config'] as Map<String, dynamic>,
      status: json['status'] as String,
      runState: json['run_state'] as Map<String, dynamic>? ?? {},
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );

Map<String, dynamic> _$StrategyModelToJson(StrategyModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'user_id': instance.userId,
      'name': instance.name,
      'type': instance.type,
      'symbol': instance.symbol,
      'config': instance.config,
      'status': instance.status,
      'run_state': instance.runState,
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
    };
