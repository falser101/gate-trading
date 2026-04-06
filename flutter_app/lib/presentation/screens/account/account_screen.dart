import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/account_provider.dart';
import '../../../data/models/balance_model.dart';
import '../../../data/models/account_detail_model.dart';

class AccountScreen extends ConsumerStatefulWidget {
  const AccountScreen({super.key});

  @override
  ConsumerState<AccountScreen> createState() => _AccountScreenState();
}

class _AccountScreenState extends ConsumerState<AccountScreen> {
  @override
  void initState() {
    super.initState();
    Future.delayed(Duration.zero, () {
      if (mounted) {
        ref.read(accountProvider.notifier).loadBalance();
        ref.read(accountProvider.notifier).loadAccountDetail();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(accountProvider);
    final accountDetail = state.accountDetail;
    final totalBalance = state.totalBalance;

    return Scaffold(
      appBar: AppBar(
        title: const Text('账户'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: state.isLoading
                ? null
                : () => ref.read(accountProvider.notifier).refresh(),
          ),
        ],
      ),
      body: state.isLoading && totalBalance == null
          ? const Center(child: CircularProgressIndicator())
          : state.error != null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(Icons.error_outline,
                          size: 48, color: Colors.red[300]),
                      const SizedBox(height: 16),
                      Text(
                        '加载失败',
                        style: TextStyle(color: Colors.red[300], fontSize: 16),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        state.error!,
                        style: TextStyle(color: Colors.grey[600], fontSize: 12),
                        textAlign: TextAlign.center,
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: () =>
                            ref.read(accountProvider.notifier).refresh(),
                        child: const Text('重试'),
                      ),
                    ],
                  ),
                )
              : RefreshIndicator(
                  onRefresh: () => ref.read(accountProvider.notifier).refresh(),
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // 账户详情卡片
                        if (accountDetail != null)
                          _buildAccountDetailCard(accountDetail),
                        if (accountDetail != null) const SizedBox(height: 24),

                        // 总资产卡片
                        _buildTotalAssetCard(totalBalance),
                        const SizedBox(height: 24),

                        // 余额明细列表
                        const Text(
                          '资产明细',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 12),
                        _buildBalanceList(state),
                      ],
                    ),
                  ),
                ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 0,
        onTap: (index) {
          if (index == 0) Navigator.pop(context);
          if (index == 1) Navigator.pop(context);
          if (index == 2) Navigator.pop(context);
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.account_balance_wallet),
            label: '账户',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.show_chart),
            label: '行情',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.receipt_long),
            label: '订单',
          ),
        ],
      ),
    );
  }

  Widget _buildBalanceList(AccountState state) {
    final totalBalance = state.totalBalance;
    final details = totalBalance?.details;

    if (details == null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.account_balance_wallet_outlined,
                size: 64, color: Colors.grey[600]),
            const SizedBox(height: 16),
            Text(
              '暂无资产',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
          ],
        ),
      );
    }

    // 获取所有非零余额的业务线
    final nonZeroItems = <MapEntry<String, BalanceItem>>[];

    if (details.crossMargin != null && _isNonZero(details.crossMargin!)) {
      nonZeroItems.add(MapEntry('跨账户', details.crossMargin!));
    }
    if (details.spot != null && _isNonZero(details.spot!)) {
      nonZeroItems.add(MapEntry('现货', details.spot!));
    }
    if (details.finance != null && _isNonZero(details.finance!)) {
      nonZeroItems.add(MapEntry('理财', details.finance!));
    }
    if (details.margin != null && _isNonZero(details.margin!)) {
      nonZeroItems.add(MapEntry('杠杆', details.margin!));
    }
    if (details.quant != null && _isNonZero(details.quant!)) {
      nonZeroItems.add(MapEntry('量化', details.quant!));
    }
    if (details.futures != null && _isNonZero(details.futures!)) {
      nonZeroItems.add(MapEntry('合约', details.futures!));
    }
    if (details.delivery != null && _isNonZero(details.delivery!)) {
      nonZeroItems.add(MapEntry('交割', details.delivery!));
    }
    if (details.warrant != null && _isNonZero(details.warrant!)) {
      nonZeroItems.add(MapEntry('权证', details.warrant!));
    }
    if (details.cbbc != null && _isNonZero(details.cbbc!)) {
      nonZeroItems.add(MapEntry('牛熊证', details.cbbc!));
    }
    if (details.memeBox != null && _isNonZero(details.memeBox!)) {
      nonZeroItems.add(MapEntry('Meme', details.memeBox!));
    }
    if (details.options != null && _isNonZero(details.options!)) {
      nonZeroItems.add(MapEntry('期权', details.options!));
    }
    if (details.payment != null && _isNonZero(details.payment!)) {
      nonZeroItems.add(MapEntry('支付', details.payment!));
    }

    if (nonZeroItems.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.account_balance_wallet_outlined,
                size: 64, color: Colors.grey[600]),
            const SizedBox(height: 16),
            Text(
              '暂无资产',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
            const SizedBox(height: 8),
            Text(
              '充值或交易后查看余额',
              style: TextStyle(color: Colors.grey[700], fontSize: 12),
            ),
          ],
        ),
      );
    }

    return ListView.separated(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: nonZeroItems.length,
      separatorBuilder: (_, __) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final entry = nonZeroItems[index];
        return _buildBalanceTile(entry.key, entry.value);
      },
    );
  }

  bool _isNonZero(BalanceItem item) {
    final amount = double.tryParse(item.amount) ?? 0;
    final borrowed = double.tryParse(item.borrowed ?? '0') ?? 0;
    return amount > 0 || borrowed > 0;
  }

  Widget _buildTotalAssetCard(TotalBalanceModel? totalBalance) {
    final total = totalBalance?.total?.amountValue ?? 0;
    final unrealisedPnl = totalBalance?.total?.unrealisedPnlValue ?? 0;
    final borrowed = totalBalance?.total?.borrowedValue ?? 0;

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF00DC82), Color(0xFF00B86C)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '总资产估值',
            style: TextStyle(color: Colors.black54, fontSize: 14),
          ),
          const SizedBox(height: 8),
          Text(
            '\$${total.toStringAsFixed(2)}',
            style: const TextStyle(
              color: Colors.black,
              fontSize: 36,
              fontWeight: FontWeight.bold,
            ),
          ),
          if (unrealisedPnl != 0) ...[
            const SizedBox(height: 4),
            Row(
              children: [
                Icon(
                  unrealisedPnl >= 0 ? Icons.trending_up : Icons.trending_down,
                  size: 16,
                  color: unrealisedPnl >= 0 ? Colors.green : Colors.red,
                ),
                const SizedBox(width: 4),
                Text(
                  '${unrealisedPnl >= 0 ? "+" : ""}${unrealisedPnl.toStringAsFixed(2)}',
                  style: TextStyle(
                    color: unrealisedPnl >= 0 ? Colors.green : Colors.red,
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  ' (未实现盈亏)',
                  style: TextStyle(color: Colors.black54, fontSize: 12),
                ),
              ],
            ),
          ],
          if (borrowed != 0) ...[
            const SizedBox(height: 4),
            Text(
              '借入：\$${borrowed.toStringAsFixed(2)}',
              style: TextStyle(color: Colors.black54, fontSize: 12),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildAccountDetailCard(AccountDetailModel detail) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF161B22),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: const Color(0xFF00DC82).withOpacity(0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.info_outline, color: const Color(0xFF00DC82)),
              const SizedBox(width: 8),
              const Text(
                '账户信息',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF00DC82),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          _buildDetailRow('用户 ID', detail.userId?.toString() ?? '-'),
          const SizedBox(height: 8),
          _buildDetailRow('账户等级', 'VIP ${detail.tier ?? 0}'),
          const SizedBox(height: 8),
          _buildDetailRow('API Key 类型', detail.key?.modeText ?? '-'),
          if (detail.ipWhitelist != null && detail.ipWhitelist!.isNotEmpty) ...[
            const SizedBox(height: 8),
            _buildDetailRow('IP 白名单', detail.ipWhitelist!.join(', ')),
          ],
          if (detail.currencyPairs != null &&
              detail.currencyPairs!.isNotEmpty) ...[
            const SizedBox(height: 8),
            _buildDetailRow(
                '允许交易对', detail.currencyPairs!.take(5).join(', ')),
            if (detail.currencyPairs!.length > 5)
              Text(
                '等共 ${detail.currencyPairs!.length} 个',
                style: TextStyle(color: Colors.grey[600], fontSize: 12),
              ),
          ],
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 80,
          child: Text(
            label,
            style: TextStyle(color: Colors.grey[600], fontSize: 13),
          ),
        ),
        Expanded(
          child: Text(
            value,
            style: const TextStyle(fontSize: 13, color: Colors.white),
          ),
        ),
      ],
    );
  }

  Widget _buildBalanceTile(String businessLine, BalanceItem balance) {
    final available = double.tryParse(balance.amount) ?? 0;
    final borrowed = double.tryParse(balance.borrowed ?? '0') ?? 0;
    final unrealisedPnl = double.tryParse(balance.unrealisedPnl ?? '0') ?? 0;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF161B22),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              // 币种图标
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: const Color(0xFF00DC82).withOpacity(0.2),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Center(
                  child: Text(
                    balance.currency.substring(0, 1),
                    style: const TextStyle(
                      color: Color(0xFF00DC82),
                      fontWeight: FontWeight.bold,
                      fontSize: 18,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),

              // 币种和业务线信息
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      balance.currency,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      businessLine,
                      style: TextStyle(color: Colors.grey[500], fontSize: 12),
                    ),
                  ],
                ),
              ),

              // 右侧金额
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    available.toStringAsFixed(4),
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (borrowed > 0) ...[
                    const SizedBox(height: 4),
                    Text(
                      '借入：${borrowed.toStringAsFixed(4)}',
                      style: TextStyle(color: Colors.orange[300], fontSize: 12),
                    ),
                  ],
                  if (unrealisedPnl != 0) ...[
                    const SizedBox(height: 2),
                    Text(
                      '${unrealisedPnl >= 0 ? "+" : ""}${unrealisedPnl.toStringAsFixed(4)}',
                      style: TextStyle(
                        color: unrealisedPnl >= 0 ? Colors.green : Colors.red,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}
