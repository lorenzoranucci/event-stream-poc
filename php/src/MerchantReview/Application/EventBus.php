<?php

declare(strict_types=1);

namespace App\MerchantReview\Application;

interface EventBus
{
    public function dispatch(IntegrationEvent $event): void;
}
