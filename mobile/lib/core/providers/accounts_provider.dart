import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/account.dart';
import 'app_providers.dart';

// Accounts state class
class AccountsState {
  final List<Account> accounts;
  final bool isLoading;
  final String? error;

  const AccountsState({
    this.accounts = const [],
    this.isLoading = false,
    this.error,
  });

  AccountsState copyWith({
    List<Account>? accounts,
    bool? isLoading,
    String? error,
  }) {
    return AccountsState(
      accounts: accounts ?? this.accounts,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

// Accounts state notifier
class AccountsNotifier extends StateNotifier<AccountsState> {
  final Ref _ref;

  AccountsNotifier(this._ref) : super(const AccountsState());

  Future<void> loadAccounts() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final providerService = _ref.read(providerServiceProvider);
      final accounts = await providerService.listAccounts();

      state = AccountsState(
        accounts: accounts,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> connectProvider({
    required String provider,
    required String token,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final providerService = _ref.read(providerServiceProvider);
      final account = await providerService.connectProvider(
        provider: provider,
        token: token,
      );

      // Add the new account to the list
      state = AccountsState(
        accounts: [...state.accounts, account],
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> disconnectAccount(String accountId) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final providerService = _ref.read(providerServiceProvider);
      await providerService.disconnectAccount(accountId);

      // Remove the account from the list
      state = AccountsState(
        accounts: state.accounts.where((a) => a.id != accountId).toList(),
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Accounts state provider
final accountsProvider = StateNotifierProvider<AccountsNotifier, AccountsState>((ref) {
  return AccountsNotifier(ref);
});
