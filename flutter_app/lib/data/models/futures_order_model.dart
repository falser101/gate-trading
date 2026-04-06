import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'futures_order_model.g.dart';

@JsonSerializable()
class FuturesOrderModel extends Equatable {
  final int? id;
  final String? contract;
  final String? size;
  final String? price;
  final String? status;
  final String? left;
  final String? fillPrice;
  final double? createTime;

  const FuturesOrderModel({
    this.id,
    this.contract,
    this.size,
    this.price,
    this.status,
    this.left,
    this.fillPrice,
    this.createTime,
  });

  factory FuturesOrderModel.fromJson(Map<String, dynamic> json) =>
      _$FuturesOrderModelFromJson(json);

  Map<String, dynamic> toJson() => _$FuturesOrderModelToJson(this);

  @override
  List<Object?> get props => [
        id,
        contract,
        size,
        price,
        status,
        left,
        fillPrice,
        createTime,
      ];

  double get sizeValue => double.tryParse(size ?? '0') ?? 0;
  double get priceValue => double.tryParse(price ?? '0') ?? 0;
  double get leftValue => double.tryParse(left ?? '0') ?? 0;
  double get fillPriceValue => double.tryParse(fillPrice ?? '0') ?? 0;

  /// 订单方向：正数为买入，负数为卖出
  int get side {
    final sizeVal = sizeValue;
    return sizeVal >= 0 ? 1 : -1;
  }

  /// 是否已成交
  bool get isFilled => status == 'finished';

  /// 是否开放中
  bool get isOpen => status == 'open';

  /// 成交时间
  DateTime get createdTime =>
      DateTime.fromMillisecondsSinceEpoch((createTime ?? 0).toInt() * 1000);
}
