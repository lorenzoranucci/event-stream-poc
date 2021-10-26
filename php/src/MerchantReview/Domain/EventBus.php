<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

interface EventBus
{
    public function dispatch(Event $event): void;
}
