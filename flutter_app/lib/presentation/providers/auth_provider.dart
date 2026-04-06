import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/network/api_client.dart';
import '../../data/repositories/auth_repository.dart';
import '../../data/models/user_model.dart';
import '../../core/router.dart';

class AuthState {
  final bool isLoading;
  final bool isAuthenticated;
  final UserModel? user;
  final String? error;

  AuthState({
    this.isLoading = false,
    this.isAuthenticated = false,
    this.user,
    this.error,
  });

  AuthState copyWith({
    bool? isLoading,
    bool? isAuthenticated,
    UserModel? user,
    String? error,
  }) {
    return AuthState(
      isLoading: isLoading ?? this.isLoading,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      user: user ?? this.user,
      error: error,
    );
  }
}

final apiClientProvider = Provider<ApiClient>((ref) => ApiClient());

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(ref.watch(apiClientProvider));
});

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authRepositoryProvider));
});

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;

  AuthNotifier(this._repository) : super(AuthState()) {
    _checkAuth();
  }

  Future<void> _checkAuth() async {
    final token = await _repository.apiClient.getToken();
    if (token != null) {
      final user = await _repository.getUser();
      if (user != null) {
        state = state.copyWith(isAuthenticated: true, user: user);
        updateAuthState(true);
      }
    }
  }

  Future<bool> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _repository.login(email, password);
      final user = await _repository.getUser();
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: user,
      );
      updateAuthState(true);
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  Future<bool> register(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _repository.register(email, password);
      final user = await _repository.getUser();
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: user,
      );
      updateAuthState(true);
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  Future<void> logout() async {
    await _repository.logout();
    state = AuthState();
    updateAuthState(false);
  }

  Future<void> bindApiKey(String apiKey, String apiSecret) async {
    await _repository.bindApiKey(apiKey, apiSecret);
    // 刷新用户数据
    final user = await _repository.getUser();
    if (user != null) {
      state = state.copyWith(user: user);
    }
  }
}
