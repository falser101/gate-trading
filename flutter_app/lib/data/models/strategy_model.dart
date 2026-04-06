import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'strategy_model.g.dart';

@JsonSerializable()
class StrategyModel extends Equatable {
  final int id;
  @JsonKey(name: 'user_id')
  final int userId;
  final String? name;
  final String type; // grid, dca
  final String symbol;
  final Map<String, dynamic> config;
  final String status; // running, stopped, paused
  @JsonKey(name: 'run_state', defaultValue: {})
  final Map<String, dynamic> runState;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const StrategyModel({
    required this.id,
    required this.userId,
    this.name,
    required this.type,
    required this.symbol,
    required this.config,
    required this.status,
    required this.runState,
    required this.createdAt,
    required this.updatedAt,
  });

  factory StrategyModel.fromJson(Map<String, dynamic> json) => _$StrategyModelFromJson(json);
  Map<String, dynamic> toJson() => _$StrategyModelToJson(this);

  String get profit => runState['total_profit']?.toString() ?? '0.00';
  int get buyTimes => runState['total_buy_times'] ?? 0;
  int get sellTimes => runState['total_sell_times'] ?? 0;

  @override
  List<Object?> get props => [id, userId, name, type, symbol, config, status, runState, createdAt, updatedAt];
}
