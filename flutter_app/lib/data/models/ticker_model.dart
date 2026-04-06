import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'ticker_model.g.dart';

@JsonSerializable()
class TickerModel extends Equatable {
  @JsonKey(name: 'currency_pair')
  final String currencyPair;
  final String last;
  @JsonKey(name: 'lowest_ask')
  final String lowestAsk;
  @JsonKey(name: 'highest_bid')
  final String highestBid;
  @JsonKey(name: 'change_percentage')
  final String changePercentage;
  @JsonKey(name: 'base_volume')
  final String baseVolume;
  @JsonKey(name: 'quote_volume')
  final String quoteVolume;
  @JsonKey(name: 'high_24h')
  final String high24h;
  @JsonKey(name: 'low_24h')
  final String low24h;

  const TickerModel({
    required this.currencyPair,
    required this.last,
    required this.lowestAsk,
    required this.highestBid,
    required this.changePercentage,
    required this.baseVolume,
    required this.quoteVolume,
    required this.high24h,
    required this.low24h,
  });

  factory TickerModel.fromJson(Map<String, dynamic> json) => _$TickerModelFromJson(json);
  Map<String, dynamic> toJson() => _$TickerModelToJson(this);

  bool get isUp {
    final change = double.tryParse(changePercentage) ?? 0;
    return change >= 0;
  }

  @override
  List<Object?> get props => [
        currencyPair,
        last,
        lowestAsk,
        highestBid,
        changePercentage,
        baseVolume,
        quoteVolume,
        high24h,
        low24h,
      ];
}
