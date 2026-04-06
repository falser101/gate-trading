// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'ticker_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

TickerModel _$TickerModelFromJson(Map<String, dynamic> json) => TickerModel(
  currencyPair: json['currency_pair'] as String,
  last: json['last'] as String,
  lowestAsk: json['lowest_ask'] as String,
  highestBid: json['highest_bid'] as String,
  changePercentage: json['change_percentage'] as String,
  baseVolume: json['base_volume'] as String,
  quoteVolume: json['quote_volume'] as String,
  high24h: json['high_24h'] as String,
  low24h: json['low_24h'] as String,
);

Map<String, dynamic> _$TickerModelToJson(TickerModel instance) =>
    <String, dynamic>{
      'currency_pair': instance.currencyPair,
      'last': instance.last,
      'lowest_ask': instance.lowestAsk,
      'highest_bid': instance.highestBid,
      'change_percentage': instance.changePercentage,
      'base_volume': instance.baseVolume,
      'quote_volume': instance.quoteVolume,
      'high_24h': instance.high24h,
      'low_24h': instance.low24h,
    };
