class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api';
  static const String healthEndpoint = '/health';

  // Auth
  static const String register = '/auth/register';
  static const String login = '/auth/login';

  // User
  static const String user = '/user';
  static const String apiKey = '/user/api-key';

  // Strategies
  static const String strategies = '/strategies';

  // Orders
  static const String orders = '/orders';

  // Market
  static const String marketTicker = '/market/ticker';

  // Account
  static const String accountBalance = '/account/balance';
  static const String accountDetail = '/account/detail';

  // Futures
  static const String futuresAccount = '/futures/accounts';
  static const String futuresPositions = '/futures/positions';
  static const String futuresPosition = '/futures/positions';
  static const String futuresPositionClose = '/futures/positions/close';
  static const String futuresOrders = '/futures/orders';
  static const String futuresOrder = '/futures/orders';
  static const String futuresTickers = '/futures/tickers';
  static const String futuresTicker = '/futures/tickers';
  static const String futuresContracts = '/futures/contracts';
}

class AppConstants {
  static const String appName = 'Gate Trading';
  static const String tokenKey = 'auth_token';
  static const String userKey = 'user_data';
}
