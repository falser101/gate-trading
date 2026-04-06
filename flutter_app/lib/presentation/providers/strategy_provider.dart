import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/network/api_client.dart';
import '../../data/repositories/strategy_repository.dart';
import '../../data/models/strategy_model.dart';

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient());

final strategyRepositoryProvider = Provider<StrategyRepository>((ref) {
  return StrategyRepository(ref.watch(apiClientProvider));
});

final strategiesProvider = StateNotifierProvider<StrategiesNotifier, StrategiesState>((ref) {
  return StrategiesNotifier(ref.watch(strategyRepositoryProvider));
});

class StrategiesState {
  final bool isLoading;
  final List<StrategyModel> strategies;
  final String? error;

  StrategiesState({
    this.isLoading = false,
    this.strategies = const [],
    this.error,
  });

  StrategiesState copyWith({
    bool? isLoading,
    List<StrategyModel>? strategies,
    String? error,
  }) {
    return StrategiesState(
      isLoading: isLoading ?? this.isLoading,
      strategies: strategies ?? this.strategies,
      error: error,
    );
  }
}

class StrategiesNotifier extends StateNotifier<StrategiesState> {
  final StrategyRepository _repository;

  StrategiesNotifier(this._repository) : super(StrategiesState());

  Future<void> loadStrategies() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final strategies = await _repository.getStrategies();
      state = state.copyWith(isLoading: false, strategies: strategies);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> createStrategy({
    required String type,
    required String symbol,
    required Map<String, dynamic> config,
  }) async {
    try {
      await _repository.createStrategy(type: type, symbol: symbol, config: config);
      await loadStrategies();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<void> startStrategy(int id) async {
    await _repository.startStrategy(id);
    await loadStrategies();
  }

  Future<void> stopStrategy(int id) async {
    await _repository.stopStrategy(id);
    await loadStrategies();
  }

  Future<void> deleteStrategy(int id) async {
    await _repository.deleteStrategy(id);
    await loadStrategies();
  }
}
