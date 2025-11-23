import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/providers/accounts_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';
import '../../../core/widgets/gradient_button.dart';

class AccountsScreen extends ConsumerStatefulWidget {
  const AccountsScreen({super.key});

  @override
  ConsumerState<AccountsScreen> createState() => _AccountsScreenState();
}

class _AccountsScreenState extends ConsumerState<AccountsScreen> {
  @override
  void initState() {
    super.initState();
    // Load accounts when screen initializes
    Future.microtask(() async {
      try {
        await ref.read(accountsProvider.notifier).loadAccounts();
      } catch (e) {
        // Handle error silently - will be shown in UI via error state
      }
    });
  }

  Future<void> _disconnectAccount(String accountId) async {
    try {
      await ref.read(accountsProvider.notifier).disconnectAccount(accountId);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: const Text('Account disconnected successfully'),
            backgroundColor: Colors.green,
            behavior: SnackBarBehavior.floating,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(e.toString()),
            backgroundColor: Colors.red,
            behavior: SnackBarBehavior.floating,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final accountsState = ref.watch(accountsProvider);

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              AppTheme.darkBackground,
              AppTheme.deepPurple.withValues(alpha: 0.3),
              AppTheme.darkBackground,
            ],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              // Header
              Padding(
                padding: const EdgeInsets.all(24),
                child: Row(
                  children: [
                    IconButton(
                      onPressed: () {
                        if (context.canPop()) {
                          context.pop();
                        } else {
                          context.go('/');
                        }
                      },
                      icon: const Icon(Icons.arrow_back, color: Colors.white),
                    ),
                    const SizedBox(width: 16),
                    const Text(
                      'Connected Accounts',
                      style: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                  ],
                ),
              ),

              // Accounts list
              Expanded(
                child: accountsState.isLoading
                    ? const Center(
                        child: CircularProgressIndicator(color: AppTheme.primaryPurple),
                      )
                    : accountsState.accounts.isEmpty
                        ? Center(
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Icon(
                                  Icons.lightbulb_outline,
                                  size: 64,
                                  color: Colors.white.withValues(alpha: 0.3),
                                ),
                                const SizedBox(height: 16),
                                Text(
                                  'No accounts connected',
                                  style: TextStyle(
                                    fontSize: 18,
                                    color: Colors.white.withValues(alpha: 0.7),
                                  ),
                                ),
                                const SizedBox(height: 32),
                                GradientButton(
                                  onPressed: () => context.push('/providers/connect'),
                                  text: 'Connect Account',
                                ),
                              ],
                            ),
                          )
                        : ListView.builder(
                            padding: const EdgeInsets.symmetric(horizontal: 24),
                            itemCount: accountsState.accounts.length,
                            itemBuilder: (context, index) {
                              final account = accountsState.accounts[index];
                              return Padding(
                                padding: const EdgeInsets.only(bottom: 16),
                                child: GlassContainer(
                                  child: ListTile(
                                    contentPadding: const EdgeInsets.all(16),
                                    leading: CircleAvatar(
                                      backgroundColor: AppTheme.primaryPurple.withValues(alpha: 0.2),
                                      child: const Icon(
                                        Icons.lightbulb,
                                        color: AppTheme.primaryPurple,
                                      ),
                                    ),
                                    title: Text(
                                      account.displayName,
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                    subtitle: Text(
                                      account.provider.toUpperCase(),
                                      style: TextStyle(
                                        color: Colors.white.withValues(alpha: 0.7),
                                      ),
                                    ),
                                    trailing: IconButton(
                                      onPressed: () => _disconnectAccount(account.id),
                                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                                    ),
                                  ),
                                ),
                              );
                            },
                          ),
              ),

              // Add account button
              if (!accountsState.isLoading && accountsState.accounts.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.all(24),
                  child: GradientButton(
                    onPressed: () => context.push('/providers/connect'),
                    text: 'Add Another Account',
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
