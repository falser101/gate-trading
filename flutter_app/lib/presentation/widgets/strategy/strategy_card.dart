import 'package:flutter/material.dart';
import '../../../data/models/strategy_model.dart';

class StrategyCard extends StatelessWidget {
  final StrategyModel strategy;
  final VoidCallback onTap;
  final VoidCallback onToggle;

  const StrategyCard({
    super.key,
    required this.strategy,
    required this.onTap,
    required this.onToggle,
  });

  @override
  Widget build(BuildContext context) {
    final isRunning = strategy.status == 'running';
    final profit = double.tryParse(strategy.profit) ?? 0;
    final isProfit = profit >= 0;

    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: const Color(0xFF161B22),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isRunning ? const Color(0xFF00DC82) : Colors.grey[800]!,
            width: 1,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: _getTypeColor(strategy.type).withOpacity(0.2),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    _getTypeName(strategy.type),
                    style: TextStyle(
                      color: _getTypeColor(strategy.type),
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const Spacer(),
                Text(
                  strategy.symbol,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '盈亏',
                      style: TextStyle(color: Colors.grey, fontSize: 12),
                    ),
                    Text(
                      '${isProfit ? '+' : ''}\$${strategy.profit}',
                      style: TextStyle(
                        color: isProfit ? const Color(0xFF00DC82) : Colors.red,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    const Text(
                      '状态',
                      style: TextStyle(color: Colors.grey, fontSize: 12),
                    ),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          decoration: BoxDecoration(
                            color: isRunning
                                ? const Color(0xFF00DC82)
                                : Colors.grey,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          _getStatusText(strategy.status),
                          style: const TextStyle(fontSize: 12),
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: onToggle,
                    style: OutlinedButton.styleFrom(
                      foregroundColor: isRunning ? Colors.red : Colors.green,
                      side: BorderSide(
                        color: isRunning ? Colors.red : Colors.green,
                      ),
                    ),
                    child: Text(isRunning ? '暂停' : '启动'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _getTypeName(String type) {
    switch (type) {
      case 'grid':
        return '网格';
      case 'dca':
        return '定投';
      default:
        return type.toUpperCase();
    }
  }

  Color _getTypeColor(String type) {
    switch (type) {
      case 'grid':
        return const Color(0xFF00DC82);
      case 'dca':
        return Colors.blue;
      default:
        return Colors.grey;
    }
  }

  String _getStatusText(String status) {
    switch (status) {
      case 'running':
        return '运行中';
      case 'stopped':
        return '已停止';
      case 'paused':
        return '已暂停';
      default:
        return status;
    }
  }
}
