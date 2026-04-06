import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';
import '../models/futures_account_model.dart';
import '../models/futures_position_model.dart';
import '../models/futures_order_model.dart';

class FuturesRepository {
  final ApiClient _client;

  FuturesRepository(this._client);

  /// 获取合约账户信息
  Future<FuturesAccountModel> getAccount() async {
    final response = await _client.get(ApiConstants.futuresAccount);
    final data = response.data as Map<String, dynamic>;
    final accountData = data['data'] as Map<String, dynamic>;
    return FuturesAccountModel.fromJson(accountData);
  }

  /// 获取持仓列表
  Future<List<FuturesPositionModel>> getPositions() async {
    final response = await _client.get(ApiConstants.futuresPositions);
    final data = response.data as Map<String, dynamic>;
    final positionsData = data['data'] as List;
    return positionsData
        .map((json) => FuturesPositionModel.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// 获取单个持仓
  Future<FuturesPositionModel> getPosition(String contract) async {
    final response = await _client.get('${ApiConstants.futuresPosition}/$contract');
    final data = response.data as Map<String, dynamic>;
    final positionData = data['data'] as Map<String, dynamic>;
    return FuturesPositionModel.fromJson(positionData);
  }

  /// 创建订单
  Future<FuturesOrderModel> createOrder({
    required String contract,
    required String size,
    required String price,
    required String tif,
    bool reduceOnly = false,
    bool close = false,
    String text = '',
  }) async {
    final response = await _client.post(ApiConstants.futuresOrders, data: {
      'contract': contract,
      'size': size,
      'price': price,
      'tif': tif,
      'reduce_only': reduceOnly,
      'close': close,
      'text': text,
    });
    final data = response.data as Map<String, dynamic>;
    final orderData = data['data'] as Map<String, dynamic>;
    return FuturesOrderModel.fromJson(orderData);
  }

  /// 获取订单列表
  Future<List<FuturesOrderModel>> getOrders({
    String status = '',
    int limit = 100,
  }) async {
    final response = await _client.get(
      ApiConstants.futuresOrders,
      queryParameters: {
        if (status.isNotEmpty) 'status': status,
        'limit': limit.toString(),
      },
    );
    final data = response.data as Map<String, dynamic>;
    final ordersData = data['data'] as List;
    return ordersData
        .map((json) => FuturesOrderModel.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// 取消订单
  Future<FuturesOrderModel> cancelOrder(String orderId) async {
    final response = await _client.delete('${ApiConstants.futuresOrder}/$orderId');
    final data = response.data as Map<String, dynamic>;
    final orderData = data['data'] as Map<String, dynamic>;
    return FuturesOrderModel.fromJson(orderData);
  }

  /// 调整杠杆
  Future<Map<String, dynamic>> setLeverage({
    required String contract,
    required String leverage,
    required String marginMode,
  }) async {
    final response = await _client.post(
      '${ApiConstants.futuresPosition}/$contract/leverage',
      data: {
        'leverage': leverage,
        'margin_mode': marginMode,
      },
    );
    return response.data['data'] as Map<String, dynamic>;
  }

  /// 平仓
  Future<Map<String, dynamic>> closePosition(String contract) async {
    final response = await _client.post(
      ApiConstants.futuresPositionClose,
      data: {'contract': contract},
    );
    return response.data['data'] as Map<String, dynamic>;
  }

  /// 获取行情列表
  Future<List<Map<String, dynamic>>> getTickers() async {
    final response = await _client.get(ApiConstants.futuresTickers);
    final data = response.data as Map<String, dynamic>;
    final tickersData = data['data'] as List;
    return tickersData.cast<Map<String, dynamic>>();
  }

  /// 获取单个行情
  Future<Map<String, dynamic>> getTicker(String contract) async {
    final response = await _client.get('${ApiConstants.futuresTicker}/$contract');
    final data = response.data as Map<String, dynamic>;
    return data['data'] as Map<String, dynamic>;
  }

  /// 获取合约列表
  Future<List<Map<String, dynamic>>> getContracts() async {
    final response = await _client.get(ApiConstants.futuresContracts);
    final data = response.data as Map<String, dynamic>;
    final contractsData = data['data'] as List;
    return contractsData.cast<Map<String, dynamic>>();
  }
}
