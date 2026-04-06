import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../data/models/futures_position_model.dart';

class FuturesPositionList extends StatelessWidget {
  final List<FuturesPositionModel> positions;
  final Function(String contract) onClose;

  const FuturesPositionList({
    super.key,
    required this.positions,
    required this.onClose,
  });

  @override
  Widget build(BuildContext context) {
    if (positions.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.pie_chart_outline, size: 64, color: Colors.grey[600]),
            const SizedBox(height: 16),
            Text(
              '暂无持仓',
              style: TextStyle(color: Colors.grey[600], fontSize: 16),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      itemCount: positions.length,
      itemBuilder: (context, index) {
        final position = positions[index];
        final isLong = position.direction > 0;
        final isProfit = position.unrealisedPnlValue >= 0;

        return Container(
          margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: const Color(0xFF161B22),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 合约名称和方向
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: isLong
                          ? const Color(0xFF00DC82).withOpacity(0.2)
                          : Colors.red.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      isLong ? '多' : '空',
                      style: TextStyle(
                        color: isLong ? const Color(0xFF00DC82) : Colors.red,
                        fontWeight: FontWeight.bold,
                        fontSize: 12,
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    position.contract ?? '',
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const Spacer(),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.grey[800],
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      '${position.leverageValue.toInt()}x',
                      style: const TextStyle(fontSize: 12),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // 盈亏信息
              Row(
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '未实现盈亏',
                        style: TextStyle(color: Colors.grey[600], fontSize: 12),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${isProfit ? '+' : ''}${position.unrealisedPnlValue.toStringAsFixed(2)} USDT',
                        style: TextStyle(
                          color: isProfit ? const Color(0xFF00DC82) : Colors.red,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const Spacer(),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text(
                        '收益率',
                        style: TextStyle(color: Colors.grey[600], fontSize: 12),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${isProfit ? '+' : ''}${position.roi.toStringAsFixed(2)}%',
                        style: TextStyle(
                          color: isProfit ? const Color(0xFF00DC82) : Colors.red,
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // 持仓详情
              Row(
                children: [
                  _buildInfoItem('持仓数量', '${position.sizeValue.toStringAsFixed(4)} USDT'),
                  _buildInfoItem('保证金', '${position.marginValue.toStringAsFixed(2)} USDT'),
                  _buildInfoItem('开仓均价', position.entryPriceValue.toStringAsFixed(2)),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  _buildInfoItem('标记价格', position.markPriceValue.toStringAsFixed(2)),
                  _buildInfoItem('强平价格', position.liqPriceValue.toStringAsFixed(2)),
                ],
              ),
              const SizedBox(height: 16),

              // 操作按钮
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => _showReverseDialog(context, position),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: Colors.white,
                        side: const BorderSide(color: Color(0xFF00DC82)),
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text('反手'),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => _showTpSlDialog(context, position),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: Colors.white,
                        side: const BorderSide(color: Colors.orange),
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text('止盈/止损'),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: ElevatedButton(
                      onPressed: () => _showCloseConfirm(context, position),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.red,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text('平仓'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildInfoItem(String label, String value) {
    return Expanded(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(color: Colors.grey[600], fontSize: 11),
          ),
          const SizedBox(height: 2),
          Text(
            value,
            style: const TextStyle(fontSize: 13),
          ),
        ],
      ),
    );
  }

  void _showReverseDialog(BuildContext context, FuturesPositionModel position) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认反手'),
        content: Text('确定要对 ${position.contract} 进行反手操作吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              // TODO: 实现反手功能
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showTpSlDialog(BuildContext context, FuturesPositionModel position) {
    final tpController = TextEditingController();
    final slController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('设置止盈止损'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: tpController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: '止盈价格',
                hintText: '输入止盈价格',
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: slController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: '止损价格',
                hintText: '输入止损价格',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              // TODO: 实现止盈止损设置
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showCloseConfirm(BuildContext context, FuturesPositionModel position) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认平仓'),
        content: Text('确定要平仓 ${position.contract ?? ''} 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              onClose(position.contract ?? '');
              Navigator.pop(context);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }
}
