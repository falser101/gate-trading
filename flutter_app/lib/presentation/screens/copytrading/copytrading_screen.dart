import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../data/models/trader_model.dart';
import '../../../data/repositories/copytrading_repository.dart';
import '../../../presentation/providers/copytrading_provider.dart';

class CopytradingScreen extends StatefulWidget {
  const CopytradingScreen({super.key});

  @override
  State<CopytradingScreen> createState() => _CopytradingScreenState();
}

class _CopytradingScreenState extends State<CopytradingScreen> {
  final _repository = CopytradingRepository();
  List<Trader> _traders = [];
  bool _isLoading = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadTraders();
  }

  Future<void> _loadTraders() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      _traders = await _repository.getTraders();
    } catch (e) {
      setState(() {
        _error = e.toString();
      });
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('跟单交易'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: RefreshIndicator(
        onRefresh: _loadTraders,
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : _error != null
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.error_outline, size: 48, color: Colors.red),
                        const SizedBox(height: 16),
                        Text('加载失败：$_error'),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: _loadTraders,
                          child: const Text('重试'),
                        ),
                      ],
                    ),
                  )
                : _traders.isEmpty
                    ? const Center(child: Text('暂无交易员数据'))
                    : ListView.builder(
                        padding: const EdgeInsets.all(8),
                        itemCount: _traders.length,
                        itemBuilder: (context, index) {
                          final trader = _traders[index];
                          return _TraderCard(trader: trader);
                        },
                      ),
      ),
    );
  }
}

class _TraderCard extends StatelessWidget {
  final Trader trader;

  const _TraderCard({required this.trader});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          children: [
            // 头像
            CircleAvatar(
              radius: 24,
              backgroundImage: trader.avatar.isNotEmpty
                  ? NetworkImage(trader.avatar)
                  : null,
              child: trader.avatar.isEmpty
                  ? const Icon(Icons.person)
                  : null,
            ),
            const SizedBox(width: 12),
            // 信息
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    trader.traderName,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      _StatChip(
                        label: '收益',
                        value: '${trader.totalRoi}%',
                        color: _getProfitColor(trader.totalRoi),
                      ),
                      const SizedBox(width: 8),
                      _StatChip(
                        label: '胜率',
                        value: '${trader.winRate}%',
                        color: Colors.blue,
                      ),
                      const SizedBox(width: 8),
                      _StatChip(
                        label: '粉丝',
                        value: trader.followerCount.toString(),
                        color: Colors.grey,
                      ),
                    ],
                  ),
                ],
              ),
            ),
            // 关注按钮
            ElevatedButton(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('关注 ${trader.traderName} 功能开发中...')),
                );
              },
              child: const Text('关注'),
            ),
          ],
        ),
      ),
    );
  }

  Color _getProfitColor(String roiStr) {
    final roi = double.tryParse(roiStr) ?? 0;
    return roi >= 0 ? Colors.green : Colors.red;
  }
}

class _StatChip extends StatelessWidget {
  final String label;
  final String value;
  final Color color;

  const _StatChip({
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Column(
        children: [
          Text(
            label,
            style: TextStyle(fontSize: 10, color: color),
          ),
          Text(
            value,
            style: TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
        ],
      ),
    );
  }
}
