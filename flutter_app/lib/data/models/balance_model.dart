import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'balance_model.g.dart';

@JsonSerializable()
class TotalBalanceModel extends Equatable {
  final BalanceDetails? details;
  final BalanceTotal? total;

  const TotalBalanceModel({
    this.details,
    this.total,
  });

  factory TotalBalanceModel.fromJson(Map<String, dynamic> json) =>
      _$TotalBalanceModelFromJson(json);
  Map<String, dynamic> toJson() => _$TotalBalanceModelToJson(this);

  @override
  List<Object?> get props => [details, total];
}

@JsonSerializable()
class BalanceDetails extends Equatable {
  @JsonKey(name: 'cross_margin')
  final BalanceItem? crossMargin;
  final BalanceItem? spot;
  final BalanceItem? finance;
  final BalanceItem? margin;
  final BalanceItem? quant;
  final BalanceItem? futures;
  final BalanceItem? delivery;
  final BalanceItem? warrant;
  final BalanceItem? cbbc;
  final BalanceItem? memeBox;
  final BalanceItem? options;
  final BalanceItem? payment;

  const BalanceDetails({
    this.crossMargin,
    this.spot,
    this.finance,
    this.margin,
    this.quant,
    this.futures,
    this.delivery,
    this.warrant,
    this.cbbc,
    this.memeBox,
    this.options,
    this.payment,
  });

  factory BalanceDetails.fromJson(Map<String, dynamic> json) =>
      _$BalanceDetailsFromJson(json);
  Map<String, dynamic> toJson() => _$BalanceDetailsToJson(this);

  List<BalanceItem> get allItems {
    return [
      if (crossMargin != null) crossMargin!,
      if (spot != null) spot!,
      if (finance != null) finance!,
      if (margin != null) margin!,
      if (quant != null) quant!,
      if (futures != null) futures!,
      if (delivery != null) delivery!,
      if (warrant != null) warrant!,
      if (cbbc != null) cbbc!,
      if (memeBox != null) memeBox!,
      if (options != null) options!,
      if (payment != null) payment!,
    ];
  }

  @override
  List<Object?> get props => [
        crossMargin,
        spot,
        finance,
        margin,
        quant,
        futures,
        delivery,
        warrant,
        cbbc,
        memeBox,
        options,
        payment,
      ];
}

@JsonSerializable()
class BalanceItem extends Equatable {
  final String amount;
  final String currency;
  final String? borrowed;
  @JsonKey(name: 'unrealised_pnl')
  final String? unrealisedPnl;

  const BalanceItem({
    required this.amount,
    required this.currency,
    this.borrowed,
    this.unrealisedPnl,
  });

  factory BalanceItem.fromJson(Map<String, dynamic> json) =>
      _$BalanceItemFromJson(json);
  Map<String, dynamic> toJson() => _$BalanceItemToJson(this);

  @override
  List<Object?> get props => [amount, currency, borrowed, unrealisedPnl];

  double get amountValue => double.tryParse(amount) ?? 0;
  double get borrowedValue => double.tryParse(borrowed ?? '0') ?? 0;
  double get unrealisedPnlValue => double.tryParse(unrealisedPnl ?? '0') ?? 0;
}

@JsonSerializable()
class BalanceTotal extends Equatable {
  final String amount;
  final String currency;
  @JsonKey(name: 'unrealised_pnl')
  final String? unrealisedPnl;
  final String? borrowed;

  const BalanceTotal({
    required this.amount,
    required this.currency,
    this.unrealisedPnl,
    this.borrowed,
  });

  factory BalanceTotal.fromJson(Map<String, dynamic> json) =>
      _$BalanceTotalFromJson(json);
  Map<String, dynamic> toJson() => _$BalanceTotalToJson(this);

  @override
  List<Object?> get props => [amount, currency, unrealisedPnl, borrowed];

  double get amountValue => double.tryParse(amount) ?? 0;
  double get borrowedValue => double.tryParse(borrowed ?? '0') ?? 0;
  double get unrealisedPnlValue => double.tryParse(unrealisedPnl ?? '0') ?? 0;
}
