import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'futures_account_model.g.dart';

@JsonSerializable()
class FuturesAccountModel extends Equatable {
  final String total;
  final String unrealisedPnl;
  final String available;
  final String orderMargin;
  final String positionMargin;
  final String maintenanceMargin;
  final String currency;
  final bool inDualMode;

  const FuturesAccountModel({
    required this.total,
    required this.unrealisedPnl,
    required this.available,
    required this.orderMargin,
    required this.positionMargin,
    required this.maintenanceMargin,
    required this.currency,
    required this.inDualMode,
  });

  factory FuturesAccountModel.fromJson(Map<String, dynamic> json) =>
      _$FuturesAccountModelFromJson(json);

  Map<String, dynamic> toJson() => _$FuturesAccountModelToJson(this);

  @override
  List<Object?> get props => [
        total,
        unrealisedPnl,
        available,
        orderMargin,
        positionMargin,
        maintenanceMargin,
        currency,
        inDualMode,
      ];

  double get totalValue => double.tryParse(total) ?? 0;
  double get unrealisedPnlValue => double.tryParse(unrealisedPnl) ?? 0;
  double get availableValue => double.tryParse(available) ?? 0;
  double get orderMarginValue => double.tryParse(orderMargin) ?? 0;
  double get positionMarginValue => double.tryParse(positionMargin) ?? 0;
  double get maintenanceMarginValue => double.tryParse(maintenanceMargin) ?? 0;
}
