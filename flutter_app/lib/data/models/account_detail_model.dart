import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'account_detail_model.g.dart';

@JsonSerializable()
class AccountDetailModel extends Equatable {
  final int? userId;
  @JsonKey(name: 'ip_whitelist')
  final List<String>? ipWhitelist;
  @JsonKey(name: 'currency_pairs')
  final List<String>? currencyPairs;
  final AccountKeyInfo? key;
  final int? tier;
  @JsonKey(name: 'copy_trading_role')
  final int? copyTradingRole;

  const AccountDetailModel({
    this.userId,
    this.ipWhitelist,
    this.currencyPairs,
    this.key,
    this.tier,
    this.copyTradingRole,
  });

  factory AccountDetailModel.fromJson(Map<String, dynamic> json) =>
      _$AccountDetailModelFromJson(json);
  Map<String, dynamic> toJson() => _$AccountDetailModelToJson(this);

  @override
  List<Object?> get props => [
        userId,
        ipWhitelist,
        currencyPairs,
        key,
        tier,
        copyTradingRole,
      ];
}

@JsonSerializable()
class AccountKeyInfo extends Equatable {
  final int? mode;

  const AccountKeyInfo({this.mode});

  factory AccountKeyInfo.fromJson(Map<String, dynamic> json) =>
      _$AccountKeyInfoFromJson(json);
  Map<String, dynamic> toJson() => _$AccountKeyInfoToJson(this);

  @override
  List<Object?> get props => [mode];

  String get modeText {
    switch (mode) {
      case 1:
        return '现货';
      case 2:
        return '合约';
      case 3:
        return '全能';
      default:
        return '未知';
    }
  }
}
