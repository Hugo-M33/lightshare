import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/providers/auth_provider.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/widgets/glass_container.dart';
import '../../../core/widgets/gradient_button.dart';

class EmailVerificationScreen extends ConsumerStatefulWidget {
  final String? token;

  const EmailVerificationScreen({
    super.key,
    this.token,
  });

  @override
  ConsumerState<EmailVerificationScreen> createState() =>
      _EmailVerificationScreenState();
}

class _EmailVerificationScreenState
    extends ConsumerState<EmailVerificationScreen> {
  bool _isVerifying = false;
  bool _isVerified = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    if (widget.token != null) {
      // Delay verification until after widget tree is built
      Future.microtask(() => _verifyEmail());
    }
  }

  Future<void> _verifyEmail() async {
    if (widget.token == null) return;

    if (mounted) {
      setState(() {
        _isVerifying = true;
        _errorMessage = null;
      });
    }

    try {
      await ref.read(authProvider.notifier).verifyEmail(widget.token!);

      if (mounted) {
        setState(() {
          _isVerifying = false;
          _isVerified = true;
        });
      }

      // Auto-navigate to home after 1.5 seconds (user is now logged in)
      await Future.delayed(const Duration(milliseconds: 1500));
      if (mounted) {
        context.go('/');
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isVerifying = false;
          _errorMessage = e.toString();
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              AppTheme.darkBackground,
              AppTheme.primaryPurple.withValues(alpha: 0.2),
              AppTheme.darkBackground,
            ],
          ),
        ),
        child: SafeArea(
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  _buildIcon(),
                  const SizedBox(height: 48),
                  GlassContainer(
                    padding: const EdgeInsets.all(32),
                    child: Column(
                      children: [
                        if (_isVerifying) _buildVerifying(),
                        if (_isVerified) _buildSuccess(),
                        if (_errorMessage != null) _buildError(),
                        if (!_isVerifying &&
                            !_isVerified &&
                            _errorMessage == null)
                          _buildNoToken(),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildIcon() {
    IconData icon;
    Color color;

    if (_isVerifying) {
      icon = Icons.hourglass_empty;
      color = AppTheme.primaryPurple;
    } else if (_isVerified) {
      icon = Icons.check_circle_outline;
      color = Colors.green;
    } else if (_errorMessage != null) {
      icon = Icons.error_outline;
      color = Colors.red;
    } else {
      icon = Icons.mark_email_unread_outlined;
      color = AppTheme.accentPink;
    }

    return Container(
      padding: const EdgeInsets.all(32),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          colors: [
            color.withValues(alpha: 0.3),
            color.withValues(alpha: 0.1),
          ],
        ),
        boxShadow: [
          BoxShadow(
            color: color.withValues(alpha: 0.5),
            blurRadius: 40,
            spreadRadius: 10,
          ),
        ],
      ),
      child: Icon(
        icon,
        size: 80,
        color: color,
      ),
    );
  }

  Widget _buildVerifying() {
    return Column(
      children: [
        const CircularProgressIndicator(
          valueColor: AlwaysStoppedAnimation(AppTheme.primaryPurple),
        ),
        const SizedBox(height: 24),
        Text(
          'Verifying your email...',
          style: Theme.of(context).textTheme.displaySmall,
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 12),
        Text(
          'Please wait while we verify your email address',
          style: Theme.of(context).textTheme.bodyMedium,
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildSuccess() {
    return Column(
      children: [
        Text(
          'Welcome to LightShare!',
          style: Theme.of(context).textTheme.displaySmall?.copyWith(
                color: Colors.green,
              ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 12),
        Text(
          'Your email has been verified and you\'re now logged in. Taking you to your dashboard...',
          style: Theme.of(context).textTheme.bodyMedium,
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildError() {
    return Column(
      children: [
        Text(
          'Verification Failed',
          style: Theme.of(context).textTheme.displaySmall?.copyWith(
                color: Colors.red,
              ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 12),
        Text(
          _errorMessage!,
          style: Theme.of(context).textTheme.bodyMedium,
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 24),
        GradientButton(
          text: 'Go to Login',
          onPressed: () {
            context.go('/auth/login');
          },
          gradientColors: const [Colors.red, Colors.orange],
        ),
      ],
    );
  }

  Widget _buildNoToken() {
    return Column(
      children: [
        Text(
          'Check Your Email',
          style: Theme.of(context).textTheme.displaySmall,
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 12),
        Text(
          'We\'ve sent a verification link to your email address. Please click the link to verify your account.',
          style: Theme.of(context).textTheme.bodyMedium,
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 24),
        GradientButton(
          text: 'Back to Login',
          onPressed: () {
            context.go('/auth/login');
          },
        ),
      ],
    );
  }
}
