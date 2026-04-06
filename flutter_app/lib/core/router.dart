import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../presentation/screens/splash/splash_screen.dart';
import '../presentation/screens/auth/login_screen.dart';
import '../presentation/screens/auth/register_screen.dart';
import '../presentation/screens/dashboard/dashboard_screen.dart';
import '../presentation/screens/strategies/grid_strategy_screen.dart';
import '../presentation/screens/strategies/dca_strategy_screen.dart';
import '../presentation/screens/market/market_screen.dart';
import '../presentation/screens/orders/orders_screen.dart';
import '../presentation/screens/settings/settings_screen.dart';
import '../presentation/screens/settings/api_key_screen.dart';
import '../presentation/screens/account/account_screen.dart';
import '../presentation/screens/futures/futures_screen.dart';
import '../data/repositories/auth_repository.dart';
import '../core/network/api_client.dart';

final _apiClient = ApiClient();
final _authNotifier = ValueNotifier<bool>(false);

/// 初始化认证状态，检查是否有 token
Future<void> initAuthState() async {
  final token = await _apiClient.getToken();
  _authNotifier.value = token != null && token.isNotEmpty;
}

/// 更新认证状态
void updateAuthState(bool isAuthenticated) {
  _authNotifier.value = isAuthenticated;
}

final GoRouter router = GoRouter(
  initialLocation: '/',
  refreshListenable: _authNotifier,
  redirect: (context, state) {
    final path = state.fullPath;
    final hasToken = _authNotifier.value;

    // Splash 页面不重定向
    if (path == '/') {
      return null;
    }

    // 认证页面：如果已登录，重定向到 dashboard
    if (path == '/login' || path == '/register') {
      if (hasToken) {
        return '/dashboard';
      }
      return null;
    }

    // 受保护的路由：如果没有 token，重定向到登录页
    if (!hasToken) {
      return '/login';
    }

    return null;
  },
  routes: [
    GoRoute(
      path: '/',
      builder: (context, state) => const SplashScreen(),
    ),
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginScreen(),
    ),
    GoRoute(
      path: '/register',
      builder: (context, state) => const RegisterScreen(),
    ),
    GoRoute(
      path: '/dashboard',
      builder: (context, state) => const DashboardScreen(),
    ),
    GoRoute(
      path: '/strategies/create/grid',
      builder: (context, state) => const GridStrategyScreen(),
    ),
    GoRoute(
      path: '/strategies/create/dca',
      builder: (context, state) => const DcaStrategyScreen(),
    ),
    GoRoute(
      path: '/market',
      builder: (context, state) => const MarketScreen(),
    ),
    GoRoute(
      path: '/orders',
      builder: (context, state) => const OrdersScreen(),
    ),
    GoRoute(
      path: '/settings',
      builder: (context, state) => const SettingsScreen(),
    ),
    GoRoute(
      path: '/settings/api-key',
      builder: (context, state) => const ApiKeyScreen(),
    ),
    GoRoute(
      path: '/account',
      builder: (context, state) => const AccountScreen(),
    ),
    GoRoute(
      path: '/futures',
      builder: (context, state) => const FuturesScreen(),
    ),
  ],
);
