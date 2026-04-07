import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'core/router.dart';
import 'core/theme/app_theme.dart';
import 'presentation/providers/copytrading_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // 初始化认证状态
  await initAuthState();
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => CopytradingProvider()),
      ],
      child: const GateTradingApp(),
    ),
  );
}

class GateTradingApp extends StatelessWidget {
  const GateTradingApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Gate Trading',
      theme: AppTheme.darkTheme,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
