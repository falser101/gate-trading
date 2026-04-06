import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/futures_provider.dart';
import '../../../data/models/futures_account_model.dart';
import 'widgets/futures_order_form.dart';
import 'widgets/futures_position_list.dart';

class FuturesScreen extends ConsumerStatefulWidget {
  const FuturesScreen({super.key});

  @override
  ConsumerState<FuturesScreen> createState() => _FuturesScreenState();
}

class _FuturesScreenState extends ConsumerState<FuturesScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    Future.delayed(Duration.zero, () {
      if (mounted) {
        ref.read(futuresProvider.notifier).loadAccount();
        ref.read(futuresProvider.notifier).loadPositions();
        ref.read(futuresProvider.notifier).loadOrders();
      }
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(futuresProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('合约交易'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '仓位'),
            Tab(text: '委托'),
            Tab(text: '资产'),
            Tab(text: '机器人'),
          ],
        ),
      ),
      body: state.isLoading && state.account == null
          ? const Center(child: CircularProgressIndicator())
          : state.error != null
              ? _buildErrorWidget(state)
              : TabBarView(
                  controller: _tabController,
                  children: [
                    _buildPositionTab(state),
                    _buildOrdersTab(state),
                    _buildAccountTab(state),
                    _buildRobotTab(),
                  ],
                ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 1,
        onTap: (index) {
          if (index == 0) Navigator.pop(context);
          if (index == 2) Navigator.pop(context);
          if (index == 3) Navigator.pop(context);
          if (index == 4) Navigator.pop(context);
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: '首页',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.show_chart),
            label: '合约',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.swap_horiz),
            label: '交易',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.account_balance_wallet),
            label: '理财',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.folder),
            label: '资产',
          ),
        ],
      ),
    );
  }

  Widget _buildErrorWidget(FuturesState state) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.error_outline, size: 48, color: Colors.red[300]),
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
            onPressed: () => ref.read(futuresProvider.notifier).refresh(),
            child: const Text('重试'),
          ),
        ],
      ),
    );
  }

  Widget _buildPositionTab(FuturesState state) {
    return Column(
      children: [
        // 账户信息摘要
        if (state.account != null) _buildAccountSummary(state.account!),
        // 订单表单
        const FuturesOrderForm(),
        // 持仓列表
        Expanded(
          child: FuturesPositionList(
            positions: state.positions,
            onClose: (contract) =>
                ref.read(futuresProvider.notifier).closePosition(contract),
          ),
        ),
      ],
    );
  }

  Widget _buildAccountSummary(FuturesAccountModel account) {
    return Container(
      padding: const EdgeInsets.all(16),
      color: const Color(0xFF161B22),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _buildSummaryItem('账户权益', account.totalValue.toStringAsFixed(2)),
              _buildSummaryItem('未实现盈亏', account.unrealisedPnlValue.toStringAsFixed(2),
                  pnlValue: account.unrealisedPnlValue),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _buildSummaryItem('可用余额', account.availableValue.toStringAsFixed(2)),
              _buildSummaryItem('持仓保证金', account.positionMarginValue.toStringAsFixed(2)),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSummaryItem(String label, String value, {double? pnlValue = 0}) {
    final isPositive = pnlValue != null && pnlValue >= 0;
    return Column(
      children: [
        Text(
          label,
          style: TextStyle(color: Colors.grey[600], fontSize: 12),
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: pnlValue != null
                ? (isPositive ? const Color(0xFF00DC82) : Colors.red)
                : const Color(0xFF00DC82),
          ),
        ),
      ],
    );
  }

  Widget _buildOrdersTab(FuturesState state) {
    if (state.orders.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.receipt_long_outlined, size: 64, color: Colors.grey[600]),
            const SizedBox(height: 16),
            Text(
              '暂无委托',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      itemCount: state.orders.length,
      itemBuilder: (context, index) {
        final order = state.orders[index];
        return ListTile(
          title: Text(order.contract ?? '', style: const TextStyle(fontWeight: FontWeight.bold)),
          subtitle: Text('价格：${order.price ?? '0'} | 数量：${order.size ?? '0'}'),
          trailing: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                order.side >= 0 ? '买入' : '卖出',
                style: TextStyle(
                  color: order.side >= 0 ? const Color(0xFF00DC82) : Colors.red,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                order.status ?? '',
                style: TextStyle(color: Colors.grey[600], fontSize: 12),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildAccountTab(FuturesState state) {
    if (state.account == null) {
      return const Center(child: CircularProgressIndicator());
    }

    final account = state.account!;
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildAccountDetailItem('总余额', account.totalValue.toStringAsFixed(2)),
        _buildAccountDetailItem('未实现盈亏', account.unrealisedPnlValue.toStringAsFixed(2)),
        _buildAccountDetailItem('可用余额', account.availableValue.toStringAsFixed(2)),
        _buildAccountDetailItem('委托保证金', account.orderMarginValue.toStringAsFixed(2)),
        _buildAccountDetailItem('持仓保证金', account.positionMarginValue.toStringAsFixed(2)),
        _buildAccountDetailItem('维持保证金', account.maintenanceMarginValue.toStringAsFixed(2)),
        _buildAccountDetailItem('结算币种', account.currency.toUpperCase()),
        _buildAccountDetailItem('双仓模式', account.inDualMode ? '已开启' : '未开启'),
      ],
    );
  }

  Widget _buildAccountDetailItem(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: TextStyle(color: Colors.grey[600])),
          Text(value, style: const TextStyle(fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildRobotTab() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.smart_toy_outlined, size: 64, color: Colors.grey[600]),
          const SizedBox(height: 16),
          Text(
            '合约机器人',
            style: TextStyle(color: Colors.grey[600], fontSize: 16),
          ),
          const SizedBox(height: 8),
          Text(
            '敬请期待',
            style: TextStyle(color: Colors.grey[700], fontSize: 12),
          ),
        ],
      ),
    );
  }
}
