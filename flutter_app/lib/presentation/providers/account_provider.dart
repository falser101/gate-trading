import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'auth_provider.dart';
import '../../data/repositories/account_repository.dart';
import '../../data/models/balance_model.dart';
import '../../data/models/account_detail_model.dart';

final accountRepositoryProvider = Provider<AccountRepository>((ref) {
  return AccountRepository(ref.watch(apiClientProvider));
});

final accountProvider =
    StateNotifierProvider<AccountNotifier, AccountState>((ref) {
  return AccountNotifier(ref.watch(accountRepositoryProvider));
});

class AccountState {
  final bool isLoading;
  final TotalBalanceModel? totalBalance;
  final AccountDetailModel? accountDetail;
  final String? error;

  AccountState({
    this.isLoading = false,
    this.totalBalance,
    this.accountDetail,
    this.error,
  });

  AccountState copyWith({
    bool? isLoading,
    TotalBalanceModel? totalBalance,
    AccountDetailModel? accountDetail,
    String? error,
  }) {
    return AccountState(
      isLoading: isLoading ?? this.isLoading,
      totalBalance: totalBalance ?? this.totalBalance,
      accountDetail: accountDetail ?? this.accountDetail,
      error: error,
    );
  }
}

class AccountNotifier extends StateNotifier<AccountState> {
  final AccountRepository _repository;

  AccountNotifier(this._repository) : super(AccountState());

  Future<void> loadBalance() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final balance = await _repository.getBalance();
      state = state.copyWith(isLoading: false, totalBalance: balance);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadAccountDetail() async {
    try {
      final detail = await _repository.getAccountDetail();
      state = state.copyWith(accountDetail: detail);
    } catch (e) {
      //  silently ignore error for account detail
    }
  }

  Future<void> refresh() async {
    await loadBalance();
    await loadAccountDetail();
  }
}
