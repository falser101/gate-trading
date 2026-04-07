class Trader {
  final String traderId;
  final String traderName;
  final String avatar;
  final String totalPnl;
  final String totalRoi;
  final String followProfit;
  final String winRate;
  final int followerCount;
  final String maxDrawdown;
  final bool isCurated;
  final bool isPrivate;
  final List<String> styleLabels;

  Trader({
    required this.traderId,
    required this.traderName,
    required this.avatar,
    required this.totalPnl,
    required this.totalRoi,
    required this.followProfit,
    required this.winRate,
    required this.followerCount,
    required this.maxDrawdown,
    required this.isCurated,
    required this.isPrivate,
    required this.styleLabels,
  });

  factory Trader.fromJson(Map<String, dynamic> json) {
    final labels = json['style_labels'] as List?;
    return Trader(
      traderId: json['trader_id'] ?? '',
      traderName: json['trader_name'] ?? '',
      avatar: json['avatar'] ?? '',
      totalPnl: json['total_pnl'] ?? '0',
      totalRoi: json['total_roi'] ?? '0',
      followProfit: json['follow_profit'] ?? '0',
      winRate: json['win_rate'] ?? '0',
      followerCount: json['follower_count'] ?? 0,
      maxDrawdown: json['max_drawdown'] ?? '0',
      isCurated: json['is_curated'] ?? false,
      isPrivate: json['is_private'] ?? false,
      styleLabels: labels?.map((e) => e.toString()).toList() ?? [],
    );
  }
}
