import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/strategy_provider.dart';

class DcaStrategyScreen extends ConsumerStatefulWidget {
  const DcaStrategyScreen({super.key});

  @override
  ConsumerState<DcaStrategyScreen> createState() => _DcaStrategyScreenState();
}

class _DcaStrategyScreenState extends ConsumerState<DcaStrategyScreen> {
  final _formKey = GlobalKey<FormState>();
  final _symbolController = TextEditingController(text: 'BTC_USDT');
  final _investAmountController = TextEditingController();
  final _intervalController = TextEditingController(text: '60');
  final _targetPriceController = TextEditingController();
  final _maxBuyTimesController = TextEditingController(text: '10');

  @override
  void dispose() {
    _symbolController.dispose();
    _investAmountController.dispose();
    _intervalController.dispose();
    _targetPriceController.dispose();
    _maxBuyTimesController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final success = await ref.read(strategiesProvider.notifier).createStrategy(
          type: 'dca',
          symbol: _symbolController.text.trim().toUpperCase(),
          config: {
            'invest_amount': _investAmountController.text,
            'interval': int.tryParse(_intervalController.text) ?? 60,
            'target_price': _targetPriceController.text,
            'max_buy_times': int.tryParse(_maxBuyTimesController.text) ?? 10,
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
        title: const Text('创建 DCA 策略'),
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
                controller: _investAmountController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '每次投资金额',
                  hintText: '如：100',
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
              const SizedBox(height: 16),
              TextFormField(
                controller: _intervalController,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: '执行间隔',
                  hintText: '如：60',
                  border: OutlineInputBorder(),
                  suffixText: '分钟',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入执行间隔';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _targetPriceController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: const InputDecoration(
                  labelText: '目标价格（可选）',
                  hintText: '低于此价格才买入',
                  border: OutlineInputBorder(),
                  suffixText: 'USDT',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _maxBuyTimesController,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: '最大买入次数',
                  hintText: '如：10',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入最大买入次数';
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
                'DCA 策略说明',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                ),
              ),
            ],
          ),
          SizedBox(height: 12),
          Text(
            '• 定期定额投资，摊薄成本\n'
            '• 适合长期定投\n'
            '• 可设置目标价格，低价时才买入',
            style: TextStyle(color: Colors.grey, fontSize: 13),
          ),
        ],
      ),
    );
  }
}
