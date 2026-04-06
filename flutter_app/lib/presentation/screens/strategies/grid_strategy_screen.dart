import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/strategy_provider.dart';

class GridStrategyScreen extends ConsumerStatefulWidget {
  const GridStrategyScreen({super.key});

  @override
  ConsumerState<GridStrategyScreen> createState() => _GridStrategyScreenState();
}

class _GridStrategyScreenState extends ConsumerState<GridStrategyScreen> {
  final _formKey = GlobalKey<FormState>();
  final _symbolController = TextEditingController(text: 'BTC_USDT');
  final _lowerPriceController = TextEditingController();
  final _upperPriceController = TextEditingController();
  final _gridCountController = TextEditingController();
  final _investAmountController = TextEditingController();

  @override
  void dispose() {
    _symbolController.dispose();
    _lowerPriceController.dispose();
    _upperPriceController.dispose();
    _gridCountController.dispose();
    _investAmountController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final success = await ref.read(strategiesProvider.notifier).createStrategy(
          type: 'grid',
          symbol: _symbolController.text.trim().toUpperCase(),
          config: {
            'lower_price': _lowerPriceController.text,
            'upper_price': _upperPriceController.text,
            'grid_count': int.tryParse(_gridCountController.text) ?? 10,
            'invest_amount': _investAmountController.text,
          },
        );

    if (!mounted) return;

    if (success) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('策略创建成功')),
      );
      context.go('/dashboard');
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('策略创建失败')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('创建网格策略'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              _buildInfoCard(),
              const SizedBox(height: 24),
              TextFormField(
                controller: _symbolController,
                decoration: const InputDecoration(
                  labelText: '交易对',
                  hintText: '如：BTC_USDT',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入交易对';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _lowerPriceController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '价格下限',
                  hintText: '如：50000',
                  border: OutlineInputBorder(),
                  suffixText: 'USDT',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入价格下限';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _upperPriceController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '价格上限',
                  hintText: '如：70000',
                  border: OutlineInputBorder(),
                  suffixText: 'USDT',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入价格上限';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _gridCountController,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: '网格数量',
                  hintText: '如：10',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入网格数量';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _investAmountController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '投资金额',
                  hintText: '如：1000',
                  border: OutlineInputBorder(),
                  suffixText: 'USDT',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入投资金额';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: _handleSubmit,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF00DC82),
                  foregroundColor: Colors.black,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                child: const Text('创建策略', style: TextStyle(fontSize: 16)),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF161B22),
        borderRadius: BorderRadius.circular(12),
      ),
      child: const Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.info_outline, color: Color(0xFF00DC82), size: 20),
              SizedBox(width: 8),
              Text(
                '网格策略说明',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                ),
              ),
            ],
          ),
          SizedBox(height: 12),
          Text(
            '• 在价格下限和上限之间自动低吸高抛\n'
            '• 适合震荡行情\n'
            '• 网格数量越多，单个网格利润越低',
            style: TextStyle(color: Colors.grey, fontSize: 13),
          ),
        ],
      ),
    );
  }
}
