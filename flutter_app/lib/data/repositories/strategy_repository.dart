import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';
import '../models/strategy_model.dart';

class StrategyRepository {
  final ApiClient _client;

  StrategyRepository(this._client);

  Future<List<StrategyModel>> getStrategies() async {
    final response = await _client.get(ApiConstants.strategies);
    final data = response.data as Map<String, dynamic>;
    final list = data['data'] as List;
    return list.map((e) => StrategyModel.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<StrategyModel> createStrategy({
    required String type,
    required String symbol,
    required Map<String, dynamic> config,
  }) async {
    final response = await _client.post(
      ApiConstants.strategies,
      data: {
        'type': type,
        'symbol': symbol,
        'config': config,
      },
    );
    return StrategyModel.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> startStrategy(int id) async {
    await _client.post('/strategies/$id/start');
  }

  Future<void> stopStrategy(int id) async {
    await _client.post('/strategies/$id/stop');
  }

  Future<void> deleteStrategy(int id) async {
    await _client.delete('/strategies/$id');
  }
}
