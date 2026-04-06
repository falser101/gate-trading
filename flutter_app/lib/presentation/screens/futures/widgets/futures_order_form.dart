import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../providers/futures_provider.dart';

class FuturesOrderForm extends ConsumerStatefulWidget {
  const FuturesOrderForm({super.key});

  @override
  ConsumerState<FuturesOrderForm> createState() => _FuturesOrderFormState();
}

class _FuturesOrderFormState extends ConsumerState<FuturesOrderForm> {
  bool _isLong = true; // true=做多，false=做空
  String _orderType = 'limit'; // limit/market
  final _priceController = TextEditingController();
  final _amountController = TextEditingController();
  String _leverage = '10';
  String _marginMode = 'cross'; // cross/isolated

  @override
  void dispose() {
    _priceController.dispose();
    _amountController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      color: const Color(0xFF161B22),
      margin: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 合约选择
          Row(
            children: [
              Expanded(
                child: DropdownButtonFormField<String>(
                  value: 'ETH_USDT',
                  decoration: InputDecoration(
                    labelText: '合约',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  items: const [
                    DropdownMenuItem(value: 'ETH_USDT', child: Text('ETH/USDT')),
                    DropdownMenuItem(value: 'BTC_USDT', child: Text('BTC/USDT')),
                  ],
                  onChanged: (value) {},
                ),
              ),
              const SizedBox(width: 12),
              // 杠杆选择
              SizedBox(
                width: 80,
                child: DropdownButtonFormField<String>(
                  value: _leverage,
                  decoration: InputDecoration(
                    labelText: '杠杆',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  items: const [
                    DropdownMenuItem(value: '1', child: Text('1x')),
                    DropdownMenuItem(value: '5', child: Text('5x')),
                    DropdownMenuItem(value: '10', child: Text('10x')),
                    DropdownMenuItem(value: '20', child: Text('20x')),
                    DropdownMenuItem(value: '50', child: Text('50x')),
                    DropdownMenuItem(value: '100', child: Text('100x')),
                  ],
                  onChanged: (value) {
                    setState(() {
                      _leverage = value!;
                    });
                  },
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 保证金模式
          Row(
            children: [
              const Text('保证金模式：'),
              ChoiceChip(
                label: const Text('全仓'),
                selected: _marginMode == 'cross',
                onSelected: (selected) {
                  if (selected) {
                    setState(() {
                      _marginMode = 'cross';
                    });
                  }
                },
              ),
              const SizedBox(width: 8),
              ChoiceChip(
                label: const Text('逐仓'),
                selected: _marginMode == 'isolated',
                onSelected: (selected) {
                  if (selected) {
                    setState(() {
                      _marginMode = 'isolated';
                    });
                  }
                },
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 订单类型
          Row(
            children: [
              Expanded(
                child: SegmentedButton<String>(
                  segments: const [
                    ButtonSegment(value: 'limit', label: Text('限价')),
                    ButtonSegment(value: 'market', label: Text('市价')),
                  ],
                  selected: {_orderType},
                  onSelectionChanged: (Set<String> selection) {
                    setState(() {
                      _orderType = selection.first;
                      if (_orderType == 'market') {
                        _priceController.text = '0';
                      }
                    });
                  },
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 价格输入
          if (_orderType == 'limit')
            TextField(
              controller: _priceController,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(
                labelText: '价格',
                suffixText: 'USDT',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
          if (_orderType == 'limit') const SizedBox(height: 12),

          // 数量输入
          TextField(
            controller: _amountController,
            keyboardType: TextInputType.number,
            decoration: InputDecoration(
              labelText: '数量',
              suffixText: '合约',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
          const SizedBox(height: 16),

          // 开多/开空按钮
          Row(
            children: [
              Expanded(
                child: ElevatedButton(
                  onPressed: () => _submitOrder(true),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF00DC82),
                    foregroundColor: Colors.black,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('开多', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: ElevatedButton(
                  onPressed: () => _submitOrder(false),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.red,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('开空', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Future<void> _submitOrder(bool isLong) async {
    final contract = 'ETH_USDT'; // TODO: 从选择器获取
    final size = _amountController.text;
    final price = _orderType == 'market' ? '0' : _priceController.text;

    if (size.isEmpty || double.tryParse(size) == null) {
      _showError('请输入有效数量');
      return;
    }

    if (_orderType == 'limit' && (price.isEmpty || double.tryParse(price) == null)) {
      _showError('请输入有效价格');
      return;
    }

    // 做多数量为正，做空数量为负
    final finalSize = isLong ? size : '-$size';

    final notifier = ref.read(futuresProvider.notifier);
    final order = await notifier.createOrder(
      contract: contract,
      size: finalSize,
      price: price,
      tif: _orderType == 'market' ? 'ioc' : 'gtc',
    );

    if (order != null) {
      _showSuccess('订单提交成功');
      _amountController.clear();
      _priceController.clear();
    } else {
      _showError('订单提交失败');
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  void _showSuccess(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: const Color(0xFF00DC82)),
    );
  }
}
