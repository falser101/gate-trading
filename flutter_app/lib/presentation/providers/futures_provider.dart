import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/auth_provider.dart';
import '../../data/repositories/futures_repository.dart';
import '../../data/models/futures_account_model.dart';
import '../../data/models/futures_position_model.dart';
import '../../data/models/futures_order_model.dart';

final futuresRepositoryProvider = Provider<FuturesRepository>((ref) {
  return FuturesRepository(ref.watch(apiClientProvider));
});

final futuresProvider =
    StateNotifierProvider<FuturesNotifier, FuturesState>((ref) {
  return FuturesNotifier(ref.watch(futuresRepositoryProvider));
});

class FuturesState {
  final bool isLoading;
  final FuturesAccountModel? account;
  final List<FuturesPositionModel> positions;
  final List<FuturesOrderModel> orders;
  final List<Map<String, dynamic>> tickers;
  final List<Map<String, dynamic>> contracts;
  final String? error;

  FuturesState({
    this.isLoading = false,
    this.account,
    this.positions = const [],
    this.orders = const [],
    this.tickers = const [],
    this.contracts = const [],
    this.error,
  });

  FuturesState copyWith({
    bool? isLoading,
    FuturesAccountModel? account,
    List<FuturesPositionModel>? positions,
    List<FuturesOrderModel>? orders,
    List<Map<String, dynamic>>? tickers,
    List<Map<String, dynamic>>? contracts,
    String? error,
  }) {
    return FuturesState(
      isLoading: isLoading ?? this.isLoading,
      account: account ?? this.account,
      positions: positions ?? this.positions,
      orders: orders ?? this.orders,
      tickers: tickers ?? this.tickers,
      contracts: contracts ?? this.contracts,
      error: error,
    );
  }
}

class FuturesNotifier extends StateNotifier<FuturesState> {
  final FuturesRepository _repository;

  FuturesNotifier(this._repository) : super(FuturesState());

  /// 加载账户信息
  Future<void> loadAccount() async {
    try {
      final account = await _repository.getAccount();
      state = state.copyWith(account: account, error: null);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载持仓
  Future<void> loadPositions() async {
    try {
      final positions = await _repository.getPositions();
      state = state.copyWith(positions: positions, error: null);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载订单
  Future<void> loadOrders({String status = '', int limit = 100}) async {
    try {
      final orders = await _repository.getOrders(status: status, limit: limit);
      state = state.copyWith(orders: orders, error: null);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 创建订单
  Future<FuturesOrderModel?> createOrder({
    required String contract,
    required String size,
    required String price,
    required String tif,
    bool reduceOnly = false,
    bool close = false,
    String text = '',
  }) async {
    try {
      final order = await _repository.createOrder(
        contract: contract,
        size: size,
        price: price,
        tif: tif,
        reduceOnly: reduceOnly,
        close: close,
        text: text,
      );
      // 刷新订单列表
      loadOrders();
      return order;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  /// 取消订单
  Future<bool> cancelOrder(String orderId) async {
    try {
      await _repository.cancelOrder(orderId);
      // 刷新订单列表
      loadOrders();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 平仓
  Future<bool> closePosition(String contract) async {
    try {
      await _repository.closePosition(contract);
      // 刷新持仓
      loadPositions();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 调整杠杆
  Future<bool> setLeverage({
    required String contract,
    required String leverage,
    required String marginMode,
  }) async {
    try {
      await _repository.setLeverage(
        contract: contract,
        leverage: leverage,
        marginMode: marginMode,
      );
      // 刷新持仓
      loadPositions();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// 加载行情
  Future<void> loadTickers() async {
    try {
      final tickers = await _repository.getTickers();
      state = state.copyWith(tickers: tickers, error: null);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 加载合约列表
  Future<void> loadContracts() async {
    try {
      final contracts = await _repository.getContracts();
      state = state.copyWith(contracts: contracts, error: null);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  /// 刷新所有数据
  Future<void> refresh() async {
    await loadAccount();
    await loadPositions();
    await loadOrders();
  }

  /// 清除错误
  void clearError() {
    state = state.copyWith(error: null);
  }
}
