import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'futures_position_model.g.dart';

@JsonSerializable()
class FuturesPositionModel extends Equatable {
  final String? contract;
  final String? size;
  final String? leverage;
  final String? entryPrice;
  final String? markPrice;
  final String? liqPrice;
  final String? margin;
  final String? unrealisedPnl;
  final String? realisedPnl;
  final String? initialMargin;
  final String? maintenanceMargin;
  final int? adlRanking;
  final String? mode;

  const FuturesPositionModel({
    this.contract,
    this.size,
    this.leverage,
    this.entryPrice,
    this.markPrice,
    this.liqPrice,
    this.margin,
    this.unrealisedPnl,
    this.realisedPnl,
    this.initialMargin,
    this.maintenanceMargin,
    this.adlRanking,
    this.mode,
  });

  factory FuturesPositionModel.fromJson(Map<String, dynamic> json) =>
      _$FuturesPositionModelFromJson(json);

  Map<String, dynamic> toJson() => _$FuturesPositionModelToJson(this);

  @override
  List<Object?> get props => [
        contract,
        size,
        leverage,
        entryPrice,
        markPrice,
        liqPrice,
        margin,
        unrealisedPnl,
        realisedPnl,
        initialMargin,
        maintenanceMargin,
        adlRanking,
        mode,
      ];

  double get sizeValue => double.tryParse(size ?? '0') ?? 0;
  double get leverageValue => double.tryParse(leverage ?? '1') ?? 1;
  double get entryPriceValue => double.tryParse(entryPrice ?? '0') ?? 0;
  double get markPriceValue => double.tryParse(markPrice ?? '0') ?? 0;
  double get liqPriceValue => double.tryParse(liqPrice ?? '0') ?? 0;
  double get marginValue => double.tryParse(margin ?? '0') ?? 0;
  double get unrealisedPnlValue => double.tryParse(unrealisedPnl ?? '0') ?? 0;
  double get realisedPnlValue => double.tryParse(realisedPnl ?? '0') ?? 0;
  double get initialMarginValue => double.tryParse(initialMargin ?? '0') ?? 0;
  double get maintenanceMarginValue => double.tryParse(maintenanceMargin ?? '0') ?? 0;

  /// 持仓方向：正数为多仓，负数为空仓，0 为无持仓
  int get direction {
    final sizeVal = sizeValue;
    if (sizeVal > 0) return 1; // 多仓
    if (sizeVal < 0) return -1; // 空仓
    return 0; // 无持仓
  }

  /// 收益率
  double get roi {
    if (initialMarginValue == 0) return 0;
    return (unrealisedPnlValue / initialMarginValue) * 100;
  }
}
