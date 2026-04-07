import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import '../../providers/auth_provider.dart';
import '../../providers/strategy_provider.dart';
import '../../widgets/strategy/strategy_card.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    // 延迟加载策略，等待认证状态确认
    Future.delayed(Duration.zero, () {
      if (mounted) {
        context.read<StrategyProvider>().loadStrategies();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final authState = context.watch<AuthProvider>();
    final strategiesState = context.watch<StrategyProvider>();
    final user = authState.user;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Gate Trading'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () => context.push('/settings'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await context.read<StrategyProvider>().loadStrategies();
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 资产卡片
              _buildPortfolioCard(),
              const SizedBox(height: 24),

              // 统计数据
              _buildStatsGrid(strategiesState.strategies),
              const SizedBox(height: 24),

              // 我的策略
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text(
                    '我的策略',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  TextButton.icon(
                    icon: const Icon(Icons.add),
                    label: const Text('新建'),
                    onPressed: () => _showStrategyTypeDialog(context),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              if (strategiesState.isLoading)
                const Center(child: CircularProgressIndicator())
              else if (strategiesState.strategies.isEmpty)
                _buildEmptyState()
              else
                ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: strategiesState.strategies.length,
                  separatorBuilder: (_, __) => const SizedBox(height: 12),
                  itemBuilder: (context, index) {
                    final strategy = strategiesState.strategies[index];
                    return StrategyCard(
                      strategy: strategy,
                      onTap: () => context.push('/strategies/${strategy.id}'),
                      onToggle: () {
                        if (strategy.status == 'running') {
                          context
                              .read<StrategyProvider>()
                              .stopStrategy(strategy.id);
                        } else {
                          context
                              .read<StrategyProvider>()
                              .startStrategy(strategy.id);
                        }
                      },
                    );
                  },
                ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 0,
        onTap: (index) {
          if (index == 1) context.push('/copytrading');
          if (index == 2) context.push('/futures');
          if (index == 3) context.push('/market');
          if (index == 4) context.push('/orders');
          if (index == 5) context.push('/settings');
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.dashboard),
            label: '首页',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.people),
            label: '跟单',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.show_chart),
            label: '合约',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.receipt_long),
            label: '订单',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.settings),
            label: '设置',
          ),
        ],
      ),
    );
  }

  Widget _buildPortfolioCard() {
    return InkWell(
      onTap: () => context.push('/account'),
      child: Container(
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
              '总资产 (USDT)',
              style: TextStyle(color: Colors.black54, fontSize: 14),
            ),
            const SizedBox(height: 8),
            const Text(
              '\$0.00',
              style: TextStyle(
                color: Colors.black,
                fontSize: 36,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: Colors.black12,
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: const Text(
                    '+\$0.00 今日',
                    style: TextStyle(color: Colors.black),
                  ),
                ),
                const Spacer(),
                const Icon(Icons.chevron_right, color: Colors.black54),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatsGrid(List strategies) {
    final runningCount =
        strategies.where((s) => s.status == 'running').length;

    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 12,
      crossAxisSpacing: 12,
      childAspectRatio: 1.5,
      children: [
        _buildStatCard(
          '运行中',
          runningCount.toString(),
          Icons.play_circle,
        ),
        _buildStatCard(
          '总策略',
          strategies.length.toString(),
          Icons.category,
        ),
        _buildStatCard(
          '买入次数',
          '0',
          Icons.trending_up,
        ),
        _buildStatCard(
          '卖出次数',
          '0',
          Icons.trending_down,
        ),
      ],
    );
  }

  Widget _buildStatCard(String label, String value, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF161B22),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: const Color(0xFF00DC82), size: 24),
          const SizedBox(height: 8),
          Text(
            value,
            style: const TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            label,
            style: const TextStyle(color: Colors.grey, fontSize: 12),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyState() {
    return Container(
      padding: const EdgeInsets.all(32),
      decoration: BoxDecoration(
        color: const Color(0xFF161B22),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Icon(Icons.inbox_outlined, size: 48, color: Colors.grey[600]),
          const SizedBox(height: 16),
          Text(
            '暂无策略',
            style: TextStyle(color: Colors.grey[600], fontSize: 16),
          ),
          const SizedBox(height: 8),
          Text(
            '创建一个策略开始自动交易',
            style: TextStyle(color: Colors.grey[700], fontSize: 12),
          ),
        ],
      ),
    );
  }

  void _showStrategyTypeDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('选择策略类型'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.grid_on),
              title: const Text('网格交易'),
              subtitle: const Text('震荡行情神器'),
              onTap: () {
                Navigator.pop(context);
                context.push('/strategies/create/grid');
              },
            ),
            ListTile(
              leading: const Icon(Icons.autorenew),
              title: const Text('DCA 定投'),
              subtitle: const Text('长期定投策略'),
              onTap: () {
                Navigator.pop(context);
                context.push('/strategies/create/dca');
              },
            ),
          ],
        ),
      ),
    );
  }
}
