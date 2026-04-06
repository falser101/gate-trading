import '../../core/network/api_client.dart';
import '../models/ticker_model.dart';

class MarketRepository {
  final ApiClient _client;

  MarketRepository(this._client);

  Future<TickerModel> getTicker(String symbol) async {
    final response = await _client.get('/market/ticker/$symbol');
    final data = response.data as Map<String, dynamic>;
    return TickerModel.fromJson(data['data'] as Map<String, dynamic>);
  }
}
